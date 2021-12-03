// +k8s:deepcopy-gen=package
// +groupName=crddemo.k8s.io

package v1

// 最上方是Kubernetes进行代码生成要用的Annotation风格的注释，被称为Global Tags
// 第一行：为v1包中的所有类型定义DeepCopy方法；
// 第二行：定义这个包对应的crddemo API组的名字
