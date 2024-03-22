package service

import (
	storagev1 "k8s.io/api/storage/v1"
	"sort"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
)

// 用于封装排，过滤和分页的数据类型
type dataselector struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// DataCell 接口用于各种资源list 的类型转换，转化后可以使用data selector 的排序过滤分页
// 接口定义了 两个方法
type DataCell interface {
	// GetCreation 该方法用于获取资源的创建时间
	GetCreation() time.Time
	// GetName 改方法用于获取资源的Name
	GetName() string
}

// DataSelectQuery 定义过滤和分页的结构体
type DataSelectQuery struct {
	// 该结构体类型为过滤
	Filter *FilterQuery
	// 改结构体类型为分页排序
	Paginate *PaginateQuery
}

// FilterQuery 通过Name 进行过滤
type FilterQuery struct {
	Name string
}

// PaginateQuery 通过 limit 和page 进行分页 排序
type PaginateQuery struct {
	limit int
	page  int
}

// Len 实现自定义结构的排序，需要重写 Len，Swap，Less方法
// Len 方法重写 用于获取数组的长度
func (d *dataselector) Len() int {
	return len(d.GenericDataList)
}

// Swap 方法重写 用于数据的排序
func (d *dataselector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 方法重写 用于数据比较大小 此处用于时间的比较大小 近的时间 大于 远的时间
func (d *dataselector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

// Sort 重写 Len Swap Less 三个方法 使用sort.Sort 进行排序
func (d *dataselector) Sort() *dataselector {
	sort.Sort(d)
	return d
}

// Filter 过滤的方法 根据传输的 Name 参数进行过滤
func (d *dataselector) Filter() *dataselector {
	// 若传入的Name == nil 则返回所有的数据
	if d.DataSelect.Filter.Name == "" {
		return d
	}
	// 定义新的数组 用于存储过滤出来的数据 并返回
	filterList := make([]DataCell, 0, len(d.GenericDataList))
	//若传入的 Name != nil 则返回包含Name 的数组
	for _, value := range d.GenericDataList {
		// 定义一个匹配的标识
		matches := true
		// 获取 Name 的值
		objName := value.GetName()
		// 判断传入的 Name 是否包含在 value 的Name 中
		if !strings.Contains(objName, d.DataSelect.Filter.Name) {
			matches = false
			continue
		}
		if matches {
			// 匹配成功就放在 数组中存储
			filterList = append(filterList, value)
		}
	}
	// 过滤的数组 放入结构体中 返回
	d.GenericDataList = filterList
	return d
}

// Pagination 方法用于数组的分页，根据 limit 和 Page 的传参，取一定范围内的数据并返回
func (d *dataselector) Pagination() *dataselector {
	//定义变量获取变量值

	page := d.DataSelect.Paginate.page
	limit := d.DataSelect.Paginate.limit
	if len(d.GenericDataList) < limit {
		limit = len(d.GenericDataList)
	}
	//检验参数的合法性
	if limit <= 0 || page <= 0 {
		return d
	}
	// 分页算法
	startindex := limit * (page - 1)
	endindex := limit*page - 1
	// 判断数组长度和 分页最后一页索引的大小
	if endindex >= len(d.GenericDataList) {
		endindex = len(d.GenericDataList) - 1
	}
	// 取出在 start index 和 end index 直接的数据 并返回
	d.GenericDataList = d.GenericDataList[startindex : endindex+1]
	return d
}

// 定义 podCell 类型 实现DataCell接口 用于数据转换
// pod cell 实现了 DataCell接口的 两个方法
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

// 定义 DeploymentCell 类型 实现DataCell接口 用于数据转换
// deploymentCell 实现了 DataCell接口的 两个方法
type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}

// 定义 daemon setCell 类型 实现DataCell接口 用于数据转换
// daemon setCell 实现了 DataCell接口的 两个方法
type daemonsetCell appsv1.DaemonSet

func (d daemonsetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d daemonsetCell) GetName() string {
	return d.Name
}

// 定义 statefulCell 类型 实现 DataCell 接口 用于数据转换
type statefulsetCell appsv1.StatefulSet

func (s statefulsetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s statefulsetCell) GetName() string {
	return s.Name
}

// 定义 serviceCell 类型 实现 DataCell 接口 用于数据转换
type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}

// 定义 ingressCell 类型 实现 DataCell 接口 用于数据转换
type ingressCell networkv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}

// 定义 configmapCell 类型 实现 DataCell 接口 用于数据转换
type configmapCell corev1.ConfigMap

func (c configmapCell) GetCreation() time.Time {
	return c.CreationTimestamp.Time
}

func (c configmapCell) GetName() string {
	return c.Name
}

// 定义 secretCell 类型 实现 DataCell 接口 用于数据转换
type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s secretCell) GetName() string {
	return s.Name
}

// 定义 pvcCell 类型 实现 DataCell 接口 用于数据转换
type pvcCell corev1.PersistentVolumeClaim

func (p pvcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvcCell) GetName() string {
	return p.Name
}

// 定义 pvCell 类型 实现 DataCell 接口 用于数据转换
type pvCell corev1.PersistentVolume

func (pv pvCell) GetCreation() time.Time {
	return pv.CreationTimestamp.Time
}

func (pv pvCell) GetName() string {
	return pv.Name
}

//定义 nsCell 类型 实现 DataCell 接口 用于数据转换

type nsCell corev1.Namespace

func (n nsCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n nsCell) GetName() string {
	return n.Name
}

// 定义storageClassCell 类型 实现 DataCell 接口 用于数据转换
type storageClassCell storagev1.StorageClass

func (sc storageClassCell) GetCreation() time.Time {
	return sc.CreationTimestamp.Time
}

func (sc storageClassCell) GetName() string {
	return sc.Name
}

// 定义ingressClassCell 类型 实现 DataCell 接口 用于数据转换
type ingressClassCell networkv1.IngressClass

func (i ingressClassCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressClassCell) GetName() string {
	return i.Name
}
