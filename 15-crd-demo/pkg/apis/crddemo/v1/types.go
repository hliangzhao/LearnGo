package v1

// TODO：本文件需要定义MyDemo资源对象

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MyDemo 资源对象定义
type MyDemo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MyDemoSpec `json:"spec"`
	// 一般情况下还应当有Status字段
}

// TODO：上面三行代码生成注释的含义是：
//  第一行：为下面资源类型生成对应的Client代码
//  第二行：个API资源类型定义里，没有Status字段，因为MyDemo才是主类型，所以+genclient要写在MyDemo之上，
//      不用写在MyDemoList之上，这时要细心注意的
//  第三行：请在生成DeepCopy的时候，实现Kubernetes提供的runtime.Object接口。否则，在某些版本的Kubernetes里，
//      你的这个类型定义会出现编译错误

// MyDemoSpec 定义了Spec字段的具体内容
type MyDemoSpec struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MydemoList 资源对象的列表形式
type MydemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []MyDemo `json:"items"`
}
