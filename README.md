# Jaguer CN 分科会 - 抽選アプリ

## セットアップ

1. 環境変数をセットします

    ```sh
    export APP_NAME=jaguer-cn-lottery
    export GCLOUD_REGION=asia-northeast1
    ```

2. Cloud Spanner を用意します

    ```sh
    gcloud beta spanner instances create "${APP_NAME}" \
        --description "Jaguer Lottery" \
        --config "regional-${GCLOUD_REGION}" \
        --processing-units 100
    gcloud spanner databases create app --instance "${APP_NAME}" \
        --ddl='CREATE TABLE swags (id INT64, name STRING(100), stock INT64) PRIMARY KEY(id)'
    ```
