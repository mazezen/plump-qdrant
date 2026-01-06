package collections

import (
	"context"
	"fmt"
	"os"

	"github.com/mazezen/plum-qdrant/go-qdrant/constant"
	"github.com/qdrant/go-client/qdrant"
)

// CreateCollection
// https://api.qdrant.tech/api-reference/collections/create-collection
func CreateCollection(client *qdrant.Client) {
	err := client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName:  constant.CollectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size: 100,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s 集合创建成功...\n", constant.CollectionName)

	err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: constant.CollectionNameCopy,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size: 100,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s 集合创建成功...\n", constant.CollectionNameCopy)	
}


// UpdateCollectionParameters
// https://api.qdrant.tech/api-reference/collections/update-collection
func UpdateCollectionParameters(client *qdrant.Client) {
	var threshold = uint64(100000)
	err := client.UpdateCollection(context.Background(), &qdrant.UpdateCollection{
		CollectionName: constant.CollectionName,
		OptimizersConfig: &qdrant.OptimizersConfigDiff{
			IndexingThreshold: &threshold,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("索引阈值修改成功: %d\n", threshold)

	collectionInfo, err := client.GetCollectionInfo(context.Background(), constant.CollectionName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("状态: %s\n", collectionInfo.Status)
	fmt.Printf("索引阈值: %d\n", *collectionInfo.Config.OptimizerConfig.IndexingThreshold)
}

// GetCollectionDetails
// https://api.qdrant.tech/api-reference/collections/get-collection
func GetCollectionDetails(client *qdrant.Client) {

	collectionInfo, err := client.GetCollectionInfo(context.Background(), constant.CollectionName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("状态: %s\n", collectionInfo.Status)
	fmt.Printf("索引阈值: %d\n", *collectionInfo.Config.OptimizerConfig.IndexingThreshold)
	// fmt.Printf("collectionInfo: %v", collectionInfo)

	config := collectionInfo.GetConfig().GetParams().GetVectorsConfig()
	fmt.Printf("Vectors config: %+v\n", config)

	if params := config.GetParams(); params != nil {
		fmt.Printf("向量维度: %d\n", params.GetSize())
	}

	if mapConfig := config.GetParamsMap(); mapConfig != nil {
		for name, vp := range mapConfig.GetMap() {
			fmt.Printf("向量 %s 维度: %d\n", name, vp.GetSize())
		}
	}
}

// CheckCollectionExists
// https://api.qdrant.tech/api-reference/collections/collection-exists
func CheckCollectionExists(client *qdrant.Client) {
	ok, err := client.CollectionExists(context.Background(), constant.CollectionName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("集合: [%s] 是否存在: [%t]\n", constant.CollectionName, ok)
}

// DeleteCollectonByGiven
// https://api.qdrant.tech/api-reference/collections/delete-collection
func DeleteCollectonByGiven(client *qdrant.Client) {
	err := client.DeleteCollection(context.Background(), constant.CollectionName)
	if err != nil {
		panic(err)
	}
	

	ok, err := client.CollectionExists(context.Background(), constant.CollectionName)
	if err != nil {
		panic(err)
	}
	if ok {
		fmt.Printf("集合: [%s] 删除失败... \n", constant.CollectionName)
		return
	}
	fmt.Printf("集合: [%s] 删除成功... \n", constant.CollectionName)
}

// ListAllCollections
// https://api.qdrant.tech/api-reference/collections/get-collections
func ListAllCollections(client *qdrant.Client) {
	collections, err := client.ListCollections(context.Background())
	if err != nil {
		panic(err)
	}

	if len(collections) == 0 {
		fmt.Println("暂无结合")
		os.Exit(1)
	}
	
	for _, collection := range collections {
		fmt.Printf("集合: [%s]\n", collection)
	}
}