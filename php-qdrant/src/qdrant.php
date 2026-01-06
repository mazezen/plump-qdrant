<?php

namespace Mazezen\PhpQdrant;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;

class Qdrant {

    private static string $collection_name = "php-qdrant-test";
    private static bool $flag = false;

    function qdrant() {
        // ---------------------------------- 初始化客户端 ----------------------------------
        $client = new Client(
            [
                'base_uri' => '',
                'timeout' => 20,`
                'headers' => [
                    'Content-Type' => 'application/json',
                    'api-key' => ''
                ],
                'body' => '{}'
            ]
        );

        // ---------------------------------- 查询集合列表 ----------------------------------
        $collections_response = $client->request(
            'GET', 
            '/collections',
            [
                'json' => [
                    'vectors' => [
                        'size' => 3,  // 向量维度
                        'distance' => 'Cosine'  // 距离度量
                    ]
                ]
            ]
        );
        if ($collections_response->getStatusCode() == 200) {
            $collections = json_decode($collections_response->getBody())->result->collections;
            foreach ($collections as $collection) {
                if ($collection->name === self::$collection_name) {
                    self::$flag = true;
                }
            }
        } else {
            echo "查询所有集合请求失败: " . $collections_response->getReasonPhrase() . "\n";
        }

        // ---------------------------------- 创建集合 ----------------------------------
        if (!self::$flag) {
            echo "----- 集合: " . self::$collection_name ."不存在, 开始创建 -----\n";
            try {
                $create_collection_response = $client->request(
                'PUT',
                '/collections/' . self::$collection_name,
                );
                if ($create_collection_response->getStatusCode() == 200) {
                    echo "集合创建成功: " . $create_collection_response->getBody() . "\n";
                } else {
                    echo "创建集合请求失败: " . $collections_response->getReasonPhrase() . "\n";
                }
            } catch (RequestException $e) {
                echo "创建集合请求失败: " . $e->getMessage() . "\n";
            }
        }

        // ---------------------------------- 查询指定集合 ----------------------------------
        // https://api.qdrant.tech/api-reference/collections/get-collection
        try {
            $get_collection_response = $client->request(
                'get', 
                '/collections/' . self::$collection_name
            );
            if ($get_collection_response->getStatusCode() == 200) {
                $result = json_decode($get_collection_response->getBody())->result;
                echo "集合: " . self::$collection_name . "状态: " . $result->status . "\n";
                echo "集合: " . self::$collection_name . "优化器状态: " . $result->optimizer_status . "\n";
                echo "集合: " . self::$collection_name . "points count: " . $result->points_count . "\n";
            } else {
                echo "查询指定集合请求失败: " . $get_collection_response->getReasonPhrase() . "\n";
            }
        } catch (RequestException $e) {
             echo "查询指定集合请求失败: " . $e->getMessage() . "\n";
        }

        // ---------------------------------- 插入数据 ----------------------------------
        // https://api.qdrant.tech/api-reference/points/upsert-points
        try {
            $upsert_response = $client->request(
                'PUT',
                "/collections/" . self::$collection_name . "/points",
                [
                    'json' => [
                        'batch' => [
                            'ids' => [42],
                            'vectors' => [[0.1, 0.2, 0.8]], 
                        ]
                    ]
                ]
            );
            if ($upsert_response->getStatusCode() == 200) {
                echo "数据插入成功: " . $upsert_response->getBody() . "\n";
            } else {
                echo "插入数据请求失败: " . $get_collection_response->getReasonPhrase() . "\n";
            }
        } catch (RequestException $e) {
             echo "插入数据请求失败: " . $e->getMessage() . "\n";
        }

        // ---------------------------------- retrieve a point ----------------------------------
        // https://api.qdrant.tech/api-reference/points/get-point
        try {
            $reterivepoint_reponse = $client->request(
                "GET",
                "/collections/" . self::$collection_name . "/points/42",
            );
            if ($upsert_response->getStatusCode == 200) {
                echo "retrieve a point: " . $reterivepoint_reponse->getBody() . "\n";
            } else {
                echo "retrieve a point请求失败: " . $get_collection_response->getReasonPhrase() . "\n";
            }
        } catch (RequestException $e) {
            echo "retrieve a point请求失败: " . $e->getMessage() . "\n";
        }

        // ---------------------------------- retrieve points ----------------------------------
        // https://api.qdrant.tech/api-reference/points/get-points
        try {
            $reterive_points_response = $client->request(
                'POST',
                '/collections/' . self::$collection_name . '/points',
                [
                    'json' => [
                        'ids' => [42],
                    ]
                ]
            );
            if ($reterive_points_response -> getStatusCode() == 200) {
                echo "retrieve points: " . $reterive_points_response->getBody() . "\n";
            } else  {
                echo "retrieve points请求失败: " . $reterive_points_response->getReasonPhrase() . "\n";
            }
        } catch (RequestException $e) {
            echo "retrieve points请求失败: " . $e->getMessage() . "\n";
        }
    
    }

}


