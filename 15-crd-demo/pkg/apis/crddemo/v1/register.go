package v1

// TODO：本文件的作用是注册一个资源对象类型（MyDemo）给api-server

import (
	"github.com/hliangzhao/LearnGo/15-crd-demo/pkg/apis/crddemo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion 用于定义api：crddemo.k8s.io/v1
var SchemeGroupVersion = schema.GroupVersion{
	Group:   crddemo.GroupName,
	Version: crddemo.Version,
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// addKnownTypes 用于将MyDemo这个资源对象暴露给client，从而使得client在代码生成时可以知道MyDemo的定义
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&MyDemo{},
		&MydemoList{},
	)

	// crddemo.k8s.io/v1实际被注册代码
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
