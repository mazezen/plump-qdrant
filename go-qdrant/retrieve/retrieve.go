package retrieve

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mazezen/plum-qdrant/go-qdrant/constant"
	"github.com/qdrant/go-client/qdrant"
)

type Retrieve struct {
	c *qdrant.Client
}

func NewRetrieve(client *qdrant.Client) *Retrieve {
	return &Retrieve{
		c: client,
	}
}

// UpsertPoints
// https://api.qdrant.tech/api-reference/points/upsert-points
// 如果事先知道要过滤的字段, 可以先创建索引,再写入数据
func (r *Retrieve) UpsertPoints() {
	cn := "goqdrant-upsert-three"
	collections, err := r.c.ListCollections(context.Background())
	if err != nil {
		fmt.Printf("查询集合列表失败: %s\n", err.Error())
		os.Exit(1)
	}
	var flag bool
	for _, collection := range collections {
		if collection == cn {
			flag = true
		}
	}
	// ========================================================================================
	
	if !flag {
		if err := r.c.CreateCollection(context.Background(), &qdrant.CreateCollection{
			CollectionName: cn,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size: 3,
				Distance: qdrant.Distance_Cosine,
			}),
		}); err != nil {
			panic(err)
		}
		fmt.Printf("%s 集合创建成功...\n", constant.CollectionName)
	}
	// ========================================================================================

	collectionInfo, err := r.c.GetCollectionInfo(context.Background(), cn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("状态: %s\n", collectionInfo.Status)
	fmt.Printf("索引阈值: %d\n", *collectionInfo.Config.OptimizerConfig.IndexingThreshold)

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
	// ========================================================================================

	fmt.Println("开始创建索引...")
	if err := r.CreatePlayloadIndex(cn, "color"); err != nil {
		fmt.Printf("创建索引失败: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("等待索引更新完成：尝试查询已插入的点...")
	time.Sleep(10 * time.Second)
	// ========================================================================================

	response, err := r.c.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: cn,
		Points: []*qdrant.PointStruct{
			{
				Id: qdrant.NewIDNum(1),
				Vectors: qdrant.NewVectors(0.9, 0.1, 0.1),
				Payload: qdrant.NewValueMap(map[string]any{
					"color": "red",
				}),
			},
			{
				Id: qdrant.NewIDNum(2),
				Vectors: qdrant.NewVectors(0.1, 0.9, 0.1),
				Payload: qdrant.NewValueMap(map[string]any{
					"color": "green",
				}),
			},
			{
				Id: qdrant.NewIDNum(3),
				Vectors: qdrant.NewVectors(0.1, 0.1, 0.9),
				Payload: qdrant.NewValueMap(map[string]any{
					"color": "blue",
				}),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("数据写入返回结果: %v\n", response)
	fmt.Printf("数据写入状态: [%s]\n ", response.GetStatus())
	// ========================================================================================

	for {
		count, err := r.CountPoints(cn)
		if err != nil {
			fmt.Printf("查询point 数量失败 : %s\n", err.Error())
			os.Exit(1)
		}
		if count == 0 {
			fmt.Printf("查询point 数量 : %d\n", count)
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Printf("point count: %d\n", count)
		break
	}
}

// CountPoints
// https://api.qdrant.tech/api-reference/points/count-points
func (r *Retrieve) CountPoints (collectionName string) (uint64, error) {
	if r.c != nil {
		count, err := r.c.Count(context.Background(), &qdrant.CountPoints{
			CollectionName: collectionName,
			Filter: &qdrant.Filter{
				Must: []*qdrant.Condition{
					qdrant.NewMatch("color", "red"),
				},
			},
		})
		if err != nil {
			return 0, err
		}
		
		return count, nil
	}
	return 0, fmt.Errorf("客户端连接断开")
}

// CreatePlayloadIndex
// https://api.qdrant.tech/api-reference/indexes/create-field-index
func (r *Retrieve) CreatePlayloadIndex(collectionName, keyword string) error {
	if r.c != nil {
		updateResult, err := r.c.CreateFieldIndex(context.Background(), &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName: keyword,
			FieldType: qdrant.FieldType_FieldTypeKeyword.Enum(),
		})
		if err != nil {
			return err
		}
		fmt.Printf("创建索引返回结果: %v\n", updateResult)
		fmt.Printf("创建索引状态 : %d\n", updateResult.Status)
		return nil
	}
	return fmt.Errorf("客户端连接断开")
}


