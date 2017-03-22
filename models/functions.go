package models

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
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
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
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
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
		return
	}
	if _, err = (*o).Update(t); err != nil {
		err = errors.Wrap(err, "ModifyFunctionNoLock")
		retcode = consts.ERROR_CODE__DB__UPDATE
		return
	}
	return
}

func (t *Functions) InsertFunctionNoLock(o *orm.Ormer) (retcode int, err error) {
	Logger.Info("[%v] enter InsertFunctionNoLock.", t.Uri)
	defer Logger.Info("[%v] left InsertFunctionNoLock.", t.Uri)
	if o == nil {
		err = errors.New("param `orm.Ormer` ptr empty")
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
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

// 每个节点由两部分组成：1.非正则；2.正则
type Tree struct {
	Trees    []*Tree
	RegMatch *regexp.Regexp
	Name     string // 节点名称
	Value    int    // 节点值, 也称功能ID
}

func (t *Functions) TableName() string {
	return "functions"
}

var (
	RootMap map[string]*Tree = map[string]*Tree{
		"GET":    &Tree{Value: -1}, // -1: 头结点；0：中间节点； >0 ：叶子节点
		"POST":   &Tree{Value: -1},
		"PUT":    &Tree{Value: -1},
		"PATCH":  &Tree{Value: -1},
		"DELETE": &Tree{Value: -1},
	}
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

// 构建多叉树
func setupTree(root **Tree, levels []string, funcId int) (node **Tree) {
	var index, length int = 0, 0
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
			//  如果uri中带有如 /v1/accounts/:id(\+d)/invalid
			// 例子中带有指定的正则表达式，则有替换、分割等操作，Name=:id ;   RegMatch=MustCompile("\+d")
			for index = 0; (*root).Trees != nil && index < len((*root).Trees); index++ {
				if (*root).Trees[index].Name == levels[0] {
					break
				}
			}
			// 同一个父亲的子节点中没有找到相同的正则表达式
			if (*root).Trees == nil || index == len((*root).Trees) {
				(*root).Trees = append((*root).Trees, &Tree{
					RegMatch: regexp.MustCompile("^.+$"),
					Name:     levels[0],
				})
				length = len((*root).Trees)
				return setupTree(&((*root).Trees[length-1]), levels[1:], funcId)
			} else {
				// 正则表达式节点的递归操作, 则levels减少一级
				return setupTree(&((*root).Trees[index]), levels[1:], funcId)
			}
		} else {
			// 否则，为非正则表达式
			for index = 0; (*root).Trees != nil && index < len((*root).Trees); index++ {
				if (*root).Trees[index].Name == levels[0] {
					break
				}
			}
			if (*root).Trees == nil || index == len((*root).Trees) {
				(*root).Trees = append((*root).Trees, &Tree{
					Name: levels[0],
				})
				length = len((*root).Trees)
				return setupTree(&((*root).Trees[length-1]), levels[1:], funcId)
			} else {
				// 非正则表达式节点的递归操作, 则levels减少一级
				return setupTree(&((*root).Trees[index]), levels[1:], funcId)
			}
		}
	} else {
		// 如果是叶子节点
		for index = 0; (*root).Trees != nil && index < len((*root).Trees); index++ {
			if (*root).Trees[index].Name == levels[0] {
				(*root).Trees[index].Value = funcId
				return root
			}
		}
		// 如果是冒号开头，则为正则表达式
		if strings.HasPrefix(levels[0], ":") {
			// ::TODO  如果uri中带有如 /v1/accounts/:id(\+d)/invalid
			// 例子中带有指定的正则表达式，则有替换、分割等操作，Name=:id ;   RegMatch=MustCompile("\+d")
			(*root).Trees = append((*root).Trees, &Tree{
				RegMatch: regexp.MustCompile("^.+$"),
				Name:     levels[0],
				Value:    funcId,
			})
		} else {
			// 否则，为非正则表达式
			(*root).Trees = append((*root).Trees, &Tree{
				Name:  levels[0],
				Value: funcId,
			})
		}
	}
	return root // 叶子节点
}

// 初始化多叉树, 并把func表数据加载到内存中
func LoadMapTree() (err error) {
	var (
		tree1 *Tree = RootMap["GET"]
		tree2 *Tree = RootMap["POST"]
		tree3 *Tree = RootMap["PUT"]
		tree4 *Tree = RootMap["PATCH"]
		tree5 *Tree = RootMap["DELETE"]

		funcs []Functions
	)
	if funcs, err = getAllFunctions(); err != nil {
		Logger.Error(err.Error())
		return
	}
	fmt.Printf("len(func)=%d\n", len(funcs))

	for index := 0; funcs != nil && index < len(funcs); index++ {
		funcs[index].Uri = strings.TrimSpace(funcs[index].Uri)
		if funcs[index].Uri == "" || !strings.HasPrefix(funcs[index].Uri, "/") {
			fmt.Printf("databases funcs uri `[%s]` not irregular, funcId=%d\n", funcs[index].Uri, funcs[index].Id)
			return
		}
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
	if tree.Value > 0 {
		stack.PrintDeTree(tree.Name, "/", tree.Value)
	}
	for index := 0; tree.Trees != nil && index < len(tree.Trees); index++ {
		stack.Push(tree.Name)
		printTree(tree.Trees[index], stack)
		stack.Pop()
	}
	return
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

func recursiveMatchNode(trees []*Tree, fields []string) (funcId int, entityIdStr string) {
	var index = 0
	if trees == nil || len(trees) <= 0 || fields == nil || len(fields) <= 0 {
		return -1, ""
	}
	for index = 0; index < len(trees); index++ {
		if trees[index].RegMatch == nil && trees[index].Name == fields[0] {
			if trees[index].Value > 0 && len(fields) == 1 {
				fmt.Printf("value=%d\n", trees[index].Value)
				funcId = trees[index].Value
				return
			}
			return recursiveMatchNode(trees[index].Trees, fields[1:])
		}
		if trees[index].RegMatch != nil && trees[index].RegMatch.MatchString(fields[0]) {
			if _, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
				fmt.Printf("通配符WildCard=%s\n", trees[index].Name)
				entityIdStr = fields[0]
				if trees[index].Value > 0 && len(fields) == 1 {
					funcId = trees[index].Value
					return
				}
				fmt.Printf("value=%d, name=%s\n", trees[index].Value, trees[index].Name)
				return recursiveMatchNode(trees[index].Trees, fields[1:])
			}
		}
	}
	if index == len(trees) {
		// 没有匹配到
		return -1, ""
	}
	return
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
	fmt.Println("path: ", uri.Path)
	funcId, entityStr = recursiveMatchNode(RootMap[info.Method].Trees, strings.Split(uri.Path, "/")[1:])
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
		retcode = consts.ERROR_CODE__SOURCE_DATA__ILLEGAL
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
