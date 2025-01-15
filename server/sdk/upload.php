<?php

   $client = new \GuzzleHttp\Client([
            RequestOptions::VERIFY => false,
        ]);
            $res =  $client->request('POST', $url, [
                'multipart' => [
                    [
                        'name'     => 'file',
                        'filename' => $file["name"],
                        'size'=>$file["size"],
                        'contents' => fopen($file["tmp_name"],"r")
                    ],
                ]
            ]);
            $data = $res->getBody();
            $data = json_decode($data, true);

        return $data;
?>