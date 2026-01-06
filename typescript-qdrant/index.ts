import { QdrantClient } from "@qdrant/js-client-rest";
import { log } from "node:console";

const client = new QdrantClient({
  url: "",
  port: 6333,
  apiKey: "",
  https: true,
});

async function main() {
  const collection_name = "typescript-qdrant-test";
  let cn_has: boolean = false;

  // ----------------------- 获取所有集合 ------------------------
  const collections: any = await client.getCollections();
  console.log("所有集合: ", collections.collections);
  collections.collections.map((collection: any) => {
    if (collection.name === collection_name) {
      cn_has = true;
    }
  });
  console.log(`${collection_name} 是否存在: ${cn_has}`);

  // ----------------------- 创建集合 ---------------------------
  if (!cn_has) {
    const ok = await client.createCollection(collection_name, {
      vectors: {
        size: 4,
        distance: "Cosine",
      },
      optimizers_config: {
        default_segment_number: 2,
      },
      replication_factor: 2,
    });
    if (ok) {
      console.log(`集合: ${collection_name} 创建成功`);
    }
  }

  // ----------------------- 获取所有集合 ---------------------------
  const collections_new = await client.getCollections();
  console.log(`所有集合: `, collections_new.collections);

  // ----------------------- 查看集合信息 ---------------------------
  const collection_info: any = await client.getCollection(collection_name);
  console.log(`集合 ${collection_name} 的状态 ${collection_info.status}`);
  console.log(
    `集合 ${collection_name} 的向量维度 ${collection_info.config.params.vectors.size}`
  );

  // ----------------------- 创建索引 ---------------------------
  const updateResult: any = await client.createPayloadIndex(collection_name, {
    field_name: "city",
    field_schema: "keyword",
    wait: true,
  });
  console.log(`wait: `, updateResult.wait);
  console.log(`ordering: `, updateResult.ordering);

  const updateResult2: any = await client.createPayloadIndex(collection_name, {
    field_name: "count",
    field_schema: "integer",
    wait: true,
  });
  console.log(`wait: `, updateResult2.wait);
  console.log(`ordering:`, updateResult2.ordering);

  const updateResult3: any = await client.createPayloadIndex(collection_name, {
    field_name: "coords",
    field_schema: "geo",
    wait: true,
  });
  console.log(`wait: `, updateResult3.wait);
  console.log(`ordering: `, updateResult3.ordering);

  // ----------------------- 创建point ---------------------------
  await client.upsert(collection_name, {
    wait: true,
    points: [
      {
        id: 1,
        vector: [0.05, 0.61, 0.76, 0.74],
        payload: {
          city: "Berlin",
          country: "Germany",
          count: 1000000,
          square: 12.5,
          coords: { lat: 1.0, lon: 2.0 },
        },
      },
      {
        id: 2,
        vector: [0.19, 0.81, 0.75, 0.11],
        payload: { city: ["Berlin", "London"] },
      },
      {
        id: 3,
        vector: [0.36, 0.55, 0.47, 0.94],
        payload: { city: ["Berlin", "Moscow"] },
      },
      {
        id: 4,
        vector: [0.18, 0.01, 0.85, 0.8],
        payload: { city: ["London", "Moscow"] },
      },
      {
        id: "98a9a4b1-4ef2-46fb-8315-a97d874fe1d7",
        vector: [0.24, 0.18, 0.22, 0.44],
        payload: { count: [0] },
      },
      {
        id: "f0e09527-b096-42a8-94e9-ea94d342b925",
        vector: [0.35, 0.08, 0.11, 0.44],
      },
    ],
  });

  // ----------------------- 节点数量 ---------------------------
  const collectionInfo = await client.getCollection(collection_name);
  console.log(`number of points: `, collectionInfo.points_count);

  // ----------------------- 检索 ---------------------------
  const points = await client.retrieve(collection_name, {
    ids: [1, 2],
  });
  console.log(`points: `, points);

  // ----------------------- search ---------------------------
  const queryVector = [0.2, 0.1, 0.9, 0.7];
  const res1 = await client.search(collection_name, {
    vector: queryVector,
    limit: 3,
  });
  console.log("search result: ", res1);

  const resBatch = await client.searchBatch(collection_name, {
    searches: [
      {
        vector: queryVector,
        limit: 1,
      },
      {
        vector: queryVector,
        limit: 2,
      },
    ],
  });
  console.log("search batch result: ", resBatch);

  // ----------------------- search filter ---------------------------
  const res2 = await client.search(collection_name, {
    vector: queryVector,
    limit: 3,
    filter: {
      must: [
        {
          key: "city",
          match: {
            value: "Berlin",
          },
        },
      ],
    },
  });
  console.log(`search result with filter: `, res2);
}

main();
