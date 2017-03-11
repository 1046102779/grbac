package models

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/1046102779/grbac/common/consts"
	. "github.com/1046102779/grbac/logger"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

var (
	METHOD_TYPE_GET    = 10 // GET
	METHOD_TYPE_POST   = 20 // POST
	METHOD_TYPE_PUT    = 30 // PUT
	METHOD_TYPE_PATCH  = 40 // PATCH
	METHOD_TYPE_DELETE = 50 // DELETE
)

type Functions struct {
	Id                 int       `orm:"column(function_id);auto"`
	Name               string    `orm:"column(name);size(60);null"`
	RegionId           int       `orm:"column(region_id);null"`
	Uri                string    `orm:"column(uri);size(500);null"`
	ThirdRegionMarkKey string    `orm:"column(third_region_mark_key);size(300);null"`
	MethodType         int16     `orm:"column(method_type);null"`
	Status             int16     `orm:"column(status);null"`
	UpdatedAt          time.Time `orm:"column(updated_at);type(datetime);null"`
	CreatedAt          time.Time `orm:"column(created_at);type(datetime);null"`
}

func (t *Functions) ReadFunctionNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter ReadFunctionNoLock.", t.Id)
	defer Logger.Info("[%v] left ReadFunctionNoLock.", t.Id)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__PARAM__ILLEGAL
	}
	if err = (*o).Read(t); err != nil {
		err = errors.Wrap(err, "ReadFunctionNoLock")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}

func (t *Functions) UpdateFunctionNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter ModifyFunctionNoLock.", t.Id)
	defer Logger.Info("[%v] left ModifyFunctionNoLock.", t.Id)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__PARAM__ILLEGAL
		return
	}
	if _, err = (*o).Update(t); err != nil {
		err = errors.Wrap(err, "ModifyFunctionNoLock")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	return
}

func (t *Functions) InsertFunctionNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter InsertFunctionNoLock.", t.Uri)
	defer Logger.Info("[%v] left InsertFunctionNoLock.", t.Uri)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__PARAM__ILLEGAL
		return
	}
	if _, err = (*o).Insert(t); err != nil {
		err = errors.Wrap(err, "InsertFunctionNoLock")
		retcode = consts.ERROR_CODE__DB__INSERT
		return
	}
	return
}

/*
							root
						/			 \
					 GET			 POST
					/					 \
				  v1					 v2
				/	 \					/	\
		 storages	accounts		...		...
		  /		\		...
	:id(正则)  	warehouses(非正则)
	  /   \			\
func_id=12 invalid	......
			  \
			  func_id=20


 * 优先匹配非正则表达式
 * 正则次优匹配
 * 存储结构为Map存储
*/

// 正则
type RegTree struct {
	Name     string
	RegMatch *regexp.Regexp
	Tree     *Tree
	Value    int
}

type Leaf struct {
	LeafFixMap    map[string]int // 非正则
	LeafWildCards []*RegTree     // 正则
}

// 每个节点由三部分组成：1.非正则；2.正则；3.叶子
type Tree struct {
	FixMap    map[string]*Tree // 非正则
	WildCards []*RegTree       // 正则
	Leaf      []*Leaf          // 叶子
}

func (t *Functions) TableName() string {
	return "functions"
}

var (
	RootMap map[string]*Tree
)

// 获取构建树的数据源
func getAllFunctions() (funcs []Functions, err error) {
	Logger.Info("enter getAllFunctions.")
	defer Logger.Info("left getAllFunctions.")
	funcs = []Functions{}
	o := orm.NewOrm()
	_, err = o.QueryTable((&Functions{}).TableName()).Filter("status", consts.STATUS_VALID).All(&funcs)
	if err != nil {
		err = errors.Wrap(err, "getAllFunctions")
		return
	}
	return
}

// 构建树
func setupTree(root **Tree, levels []string, funcId int) (node **Tree) {
	// 根节点为空，新建节点
	if *root == nil {
		*root = new(Tree)
	}
	// levels已遍历完成，则返回
	if levels == nil && len(levels) <= 0 {
		return root
	}
	// 如果是非叶子节点
	if len(levels) > 1 {
		// 如果是冒号开头，则为正则表达式
		if strings.HasPrefix(levels[0], ":") {
			// ::TODO  如果uri中带有如 /v1/accounts/:id(\+d)/invalid
			// 例子中带有指定的正则表达式，则有替换、分割等操作，Name=:id ;   RegMatch=MustCompile("\+d")
			if (*root).WildCards == nil {
				(*root).WildCards = append((*root).WildCards,
					&RegTree{
						Name:     levels[0],
						RegMatch: regexp.MustCompile("^.+$"),
						Tree:     new(Tree),
					})
			} else {
				index := 0
				for index = 0; index < len((*root).WildCards); index++ {
					if (*root).WildCards[len((*root).WildCards)-1].Name != levels[0] {
						break
					}
				}
				if index != len((*root).WildCards) {
					(*root).WildCards = append((*root).WildCards,
						&RegTree{
							Name:     levels[0],
							RegMatch: regexp.MustCompile("^.+$"),
							Tree:     new(Tree),
						})
				}
			}
			len := len((*root).WildCards)
			// 正则表达式节点的递归操作, 则levels减少一级
			return setupTree(&((*root).WildCards[len-1].Tree), levels[1:], funcId)
		} else {
			// 否则，为非正则表达式
			if (*root).FixMap == nil {
				(*root).FixMap = map[string]*Tree{
					levels[0]: new(Tree),
				}
			} else if _, ok := (*root).FixMap[levels[0]]; !ok {
				(*root).FixMap[levels[0]] = new(Tree)
			}
			tempTree := (*root).FixMap[levels[0]]
			// 非正则表达式节点的递归操作, 则levels减少一级
			return setupTree(&tempTree, levels[1:], funcId)
		}
	} else {
		// 如果是叶子节点
		// 如果是冒号开头，则为正则表达式
		(*root).Leaf = append((*root).Leaf, new(Leaf))
		len := len((*root).Leaf)
		if strings.HasPrefix(levels[0], ":") {
			// ::TODO  如果uri中带有如 /v1/accounts/:id(\+d)/invalid
			// 例子中带有指定的正则表达式，则有替换、分割等操作，Name=:id ;   RegMatch=MustCompile("\+d")
			(*root).Leaf[len-1].LeafWildCards = append((*root).Leaf[len-1].LeafWildCards,
				&RegTree{
					Name:     levels[0],
					RegMatch: regexp.MustCompile("^.+$"),
					Tree:     new(Tree),
					Value:    funcId,
				})
		} else {
			// 否则，为非正则表达式
			(*root).Leaf[len-1].LeafFixMap = map[string]int{
				levels[0]: funcId,
			}
		}
	}
	return root // 叶子节点
}

// 初始化功能树
func LoadMapTree() (err error) {
	var (
		funcs []Functions
	)
	RootMap = map[string]*Tree{}
	funcs, err = getAllFunctions()
	if err != nil {
		Logger.Error(err.Error())
		return
	}
	RootMap["GET"] = &Tree{}
	RootMap["POST"] = &Tree{}
	RootMap["PUT"] = &Tree{}
	RootMap["PATCH"] = &Tree{}
	RootMap["DELETE"] = &Tree{}
	var (
		tree1 *Tree = RootMap["GET"]
		tree2 *Tree = RootMap["POST"]
		tree3 *Tree = RootMap["PUT"]
		tree4 *Tree = RootMap["PATCH"]
		tree5 *Tree = RootMap["DELETE"]
	)
	for index := 0; funcs != nil && index < len(funcs); index++ {
		switch int(funcs[index].MethodType) {
		case METHOD_TYPE_GET:
			setupTree(&tree1, strings.Split(funcs[index].Uri, "/")[1:], funcs[index].Id)
		case METHOD_TYPE_POST:
			setupTree(&tree2, strings.Split(funcs[index].Uri, "/")[1:], funcs[index].Id)
		case METHOD_TYPE_PUT:
			setupTree(&tree3, strings.Split(funcs[index].Uri, "/")[1:], funcs[index].Id)
		case METHOD_TYPE_PATCH:
			setupTree(&tree4, strings.Split(funcs[index].Uri, "/")[1:], funcs[index].Id)
		case METHOD_TYPE_DELETE:
			setupTree(&tree5, strings.Split(funcs[index].Uri, "/")[1:], funcs[index].Id)
		}
	}
	return
}

type Stack struct {
	elems []*string
}

func (t *Stack) Push(elem string) {
	if strings.TrimSpace(elem) == "" {
		return
	}
	t.elems = append(t.elems, &elem)
	return
}

func (t *Stack) Pop() (elem string, err error) {
	if t.elems == nil || len(t.elems) <= 0 {
		err = errors.Wrap(err, "no element in stack")
		return
	}
	elem = *t.elems[len(t.elems)-1]
	t.elems = append(t.elems[:len(t.elems)-1])
	return
}

func (t *Stack) Print(sep string, funcId int) {
	if t.elems == nil || len(t.elems) <= 0 {
		return
	}
	for index := 0; t.elems != nil && index < len(t.elems); index++ {
		if index == len(t.elems)-1 {
			fmt.Printf("%s", *t.elems[index])
		} else {
			fmt.Printf("%s%s", *t.elems[index], sep)
		}
	}
	fmt.Printf("	=    %d\n", funcId)
	return
}

func (t *Stack) PrintDeTree(leaf string, sep string, funcId int) {
	t.Push(leaf)
	t.Print(sep, funcId)
	if _, err := t.Pop(); err != nil {
		Logger.Error(err.Error())
		return
	}
	return
}

// 深度遍历树
func printTree(tree *Tree, stack *Stack) {
	if tree == nil {
		return
	}
	if tree.Leaf != nil {
		for index := 0; index < len(tree.Leaf); index++ {
			for key, value := range tree.Leaf[index].LeafFixMap {
				stack.PrintDeTree(key, "/", value)
			}
			if len(tree.Leaf[index].LeafWildCards) > 0 {
				stack.PrintDeTree(tree.Leaf[index].LeafWildCards[0].Name, "/", tree.Leaf[index].LeafWildCards[0].Value)
			}
		}
	}
	if tree.FixMap != nil {
		for key, value := range tree.FixMap {
			stack.Push(key)
			printTree(value, stack)
		}
	}
	if tree.WildCards != nil && len(tree.WildCards) > 0 {
		for index := 0; index < len(tree.WildCards); index++ {
			stack.Push(tree.WildCards[index].Name)
			printTree(tree.WildCards[index].Tree, stack)
		}
	}
	stack.Pop()
}

// 打印树
func PrintTree() {
	for key, _ := range RootMap {
		stack := new(Stack)
		printTree(RootMap[key], stack)
	}
}

func init() {
	orm.RegisterModel(new(Functions))
}

type HttpRequestInfo struct {
	Body   string `json:"body"`
	Method string `json:"method"`
	Uri    string `json:"uri"`
}

func getFuncId(fields []string, method string) (funcId int, entityStr string, retcode int, err error) {
	var (
		subIndex int  = 0
		first    bool = false // 当获取到第一个实体ID后，直接把实体ID赋值给entityStr，以后再遇到不用再赋值了，因为后续的实体ID是前一个实体ID的子资源，粒度比较小，不考虑
		// 如：/v1/wechats/:id/pay/jsapi_params/:bill_id/open_id/:open_id
		// 则 :id 为唯一资源标识符，不用考虑二级、三级等子粒度资源
	)
	tree := RootMap[method]
	fmt.Println("fields: ", fields)
	for index := 0; index < len(fields); index++ {
		// 叶子节点
		if index == len(fields)-1 {
			if tree.Leaf == nil {
				// 1.匹配不到
				return -1, "", 0, nil
			}
			// 优先匹配非正则表达式
			for subIndex = 0; subIndex < len(tree.Leaf); subIndex++ {
				if _, ok := tree.Leaf[subIndex].LeafFixMap[fields[index]]; !ok {
					continue
				} else {
					break
				}
			}
			// 已匹配到
			if subIndex != len(tree.Leaf) {
				return tree.Leaf[subIndex].LeafFixMap[fields[index]], entityStr, 0, nil
			}
			// 否则，匹配正则表达式
			for subIndex = 0; subIndex < len(tree.Leaf); subIndex++ {
				for _, wildcard := range tree.Leaf[subIndex].LeafWildCards {
					// 匹配成功
					if wildcard.RegMatch.MatchString(fields[index]) {
						if entityStr == "" {
							entityStr = fields[index]
						}
						return wildcard.Value, entityStr, 0, nil
					}
				}
			}
			// 没有匹配到
			if subIndex == len(tree.Leaf) {
				return -1, "", 0, nil
			}
		} else {
			// 非叶子节点
			// 优先匹配非正则表达式
			if _, ok := tree.FixMap[fields[index]]; ok {
				// 匹配到中间节点
				tree = tree.FixMap[fields[index]]
				continue
			}
			if tree.WildCards == nil {
				// 没有匹配到
				return -1, "", 0, nil
			}
			// flag 标志位， 如果成功正则匹配，直接goto跳转到最开始的循环处
			// 否则，没有匹配到，则退出
			flag := false
			for subIndex = 0; subIndex < len(tree.WildCards); subIndex++ {
				if tree.WildCards[subIndex].RegMatch.MatchString(fields[index]) {
					if !first {
						entityStr = fields[index]
						first = true
					}
					tree = tree.WildCards[subIndex].Tree
					flag = true
					break
				}
			}
			if !flag && subIndex == len(tree.WildCards) {
				// 没有匹配到
				return -1, "", 0, nil
			}
		}
	}
	return -1, "", 0, nil
}

func GetFuncId(info *HttpRequestInfo) (funcId int, entityStr string, retcode int, err error) {
	Logger.Info("[%v %v] enter GetFuncId.", info.Method, info.Uri)
	defer Logger.Info("[%v %v] left GetFuncId.", info.Method, info.Uri)
	var (
		uri *url.URL
	)
	if uri, err = url.Parse(info.Uri); err != nil {
		err = errors.Wrap(err, "GetFuncId")
		retcode = consts.ERROR_CODE__JSON__PARSE_FAILED
		return
	}
	if uri.Path[len(uri.Path)-1] == '/' {
		uri.Path = uri.Path[0 : len(uri.Path)-1]
	}
	// 解析path，在解析树中验证一把path，看是否能取出功能ID
	if funcId, entityStr, retcode, err = getFuncId(strings.Split(uri.Path, "/")[1:], info.Method); err != nil {
		err = errors.Wrap(err, "GetFuncId")
		return
	}
	// 1.1 获取实体ID最快的做法是，在解析查找功能ID的过程中，直接把实体ID获取到

	// 1.2 说明：本来为了代码质量和可读性，获取实体ID的做法
	//      已经获取的功能ID，从数据库读取URI.PATH, 如：src = /v1/accounts/users/:id
	//		HTTP请求的URI.PATH, 如：dest = /v1/accounts/users/32
	//       src切分为：src_fields = [v1, accounts, users, :id]
	//       desc且分为：dest_fields = [v1, accounts, users, 32]
	// 所以遍历列表，当遇到src_fields[index]!=dest_fields[index]时，则拿到dest_fields[index]转化为实体ID
	// 如果转为int型的实体ID，则说明该实体是字符串，不能作为资源的唯一标识，只有ID才能作为资源唯一标识
	fmt.Printf("funcId=%d, entityStr=%s\n", funcId, entityStr)
	return
}

// 根据HTTP方法类型，获取方法名称
func GetMethodNameByType(methodType int) (methodName string) {
	Logger.Info("[%v] enter GetMethodNameByType.", methodType)
	defer Logger.Info("[%v] left GetMethodNameByType.", methodType)
	switch methodType {
	case METHOD_TYPE_GET:
		methodName = "GET"
	case METHOD_TYPE_POST:
		methodName = "POST"
	case METHOD_TYPE_PUT:
		methodName = "PUT"
	case METHOD_TYPE_PATCH:
		methodName = "PATCH"
	case METHOD_TYPE_DELETE:
		methodName = "DELETE"
	default:
		methodName = "unkown"
	}
	return
}

// 根据HTTP方法名，获取方法类型
func GetMethodTypeByName(methodName string) (methodType int) {
	Logger.Info("[%v] enter GetMethodTypeByName.", methodName)
	defer Logger.Info("[%v] left GetMethodTypeByName.", methodName)
	switch methodName {
	case "POST":
		methodType = METHOD_TYPE_POST
	case "GET":
		methodType = METHOD_TYPE_GET
	case "PUT":
		methodType = METHOD_TYPE_PUT
	case "PATCH":
		methodType = METHOD_TYPE_PATCH
	case "DELETE":
		methodType = METHOD_TYPE_DELETE
	default:
		methodType = -1
	}
	return
}

type FunctionInfos struct {
	Id     int    `json:"function_id"`
	Method string `json:"method"`
	Name   string `json:"name"`
	Uri    string `json:"uri"`
}

// 获取功能列表与搜索
// 搜索支持：名称和URI
func GetFunctions(pageIndex int64, pageSize int64, searchKey string) (funcInfos []*FunctionInfos, count int64, realCount int64, retcode int, err error) {
	Logger.Info("[%v] enter GetFunctions.", searchKey)
	defer Logger.Info("[%v] left GetFunctions.", searchKey)
	var (
		functions []*Functions
	)
	funcInfos = []*FunctionInfos{}
	cond := orm.NewCondition()
	o := orm.NewOrm()
	qs := o.QueryTable((&Functions{}).TableName()).Filter("status", consts.STATUS_VALID)
	if strings.TrimSpace(searchKey) != "" {
		qs = qs.SetCond(cond.Or("name__icontains", searchKey).Or("uri__icontains", searchKey))
	}
	count, _ = qs.Count()
	if realCount, err = qs.Limit(pageSize, pageSize*pageIndex).All(&functions); err != nil {
		err = errors.Wrap(err, "GetFunctions")
		retcode = consts.ERROR_CODE__DB__READ
		return
	}
	for index := 0; index < int(realCount); index++ {
		funcInfos = append(funcInfos, &FunctionInfos{
			Id:     functions[index].Id,
			Method: GetMethodNameByType(int(functions[index].MethodType)),
			Name:   functions[index].Name,
			Uri:    functions[index].Uri,
		})
	}
	return
}

// 解析获取thid_region_mark_key
// eg. /v1/storages/warehouses/:id/invalid -> :id
// 错误判断，如果不是以'/'开头，则返回错误
func GetThirdReginMarkKey(uri string) (markKey string, retcode int, err error) {
	if strings.TrimSpace(uri) == "" && strings.HasPrefix(uri, "/") {
		err = errors.New("param `uri` can't be empty , and prefix '/'")
		retcode = consts.ERROR_CODE__PARAM__ILLEGAL
		return
	}
	fields := strings.Split(uri, "/")
	if fields != nil && len(fields) > 0 {
		for index := 1; index < len(fields); index++ {
			if strings.HasPrefix(fields[index], ":") {
				markKey = fields[index]
				return
			}
		}
	}
	return
}
