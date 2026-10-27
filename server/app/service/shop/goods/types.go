package goods

// SpecMaps 规格映射结构
type SpecMaps struct {
	GroupMap map[string]int64 // 规格组名称到ID的映射
	ValueMap map[string]int64 // 规格值到ID的映射
}
