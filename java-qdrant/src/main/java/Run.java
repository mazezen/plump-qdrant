import static io.qdrant.client.PointIdFactory.id;
import static io.qdrant.client.VectorFactory.vector;
import static io.qdrant.client.VectorsFactory.namedVectors;

import com.google.common.util.concurrent.ListenableFuture;
import io.qdrant.client.QdrantClient;
import io.qdrant.client.QdrantGrpcClient;
import io.qdrant.client.grpc.Collections;
import io.qdrant.client.grpc.Collections.Distance;
import io.qdrant.client.grpc.Collections.VectorParams;
import io.qdrant.client.grpc.Points;
import io.qdrant.client.grpc.Points.Vector;


import java.util.List;
import java.util.HashMap;
import java.util.Arrays;

public class Run {

    private static String COLLECTION_NAME  = "java-qdrant-test";

    private static boolean flag = false;

    public static void main(String[] args) throws Exception {
        QdrantClient client = new QdrantClient(
                QdrantGrpcClient.newBuilder(
                        "",
                        6334, true)
                        .withApiKey("")
                        .build());

        // -------------------------- 获取所有的集合 --------------------------
        List<String> strings = client.listCollectionsAsync().get();
        System.out.println("所有集合: " + strings);

        for (int i = 0; i < strings.size(); i++) {
            if (strings.get(i).equals(COLLECTION_NAME)) {
                flag = true;
            }
        }

        // -------------------------- 创建集合 --------------------------
        if (!flag) {
            Collections.CollectionOperationResponse operationResponse = client.createCollectionAsync(
                    COLLECTION_NAME,
                    VectorParams.newBuilder().setDistance(Distance.Cosine).setSize(3).build()).get();
        }

        // -------------------------- 集合信息 --------------------------
        Collections.CollectionInfo collectionInfo = client.getCollectionInfoAsync(COLLECTION_NAME).get();
        System.out.println("集合状态: " + collectionInfo.getStatus());
        System.out.println("优化器状态: " + collectionInfo.getOptimizerStatus());


    }

}
