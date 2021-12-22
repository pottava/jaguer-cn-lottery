# Jaguer CN 分科会 - 抽選アプリ

## セットアップ

1. 環境変数をセットします

    ```sh
    export APP_NAME=jaguer-cn-lottery
    export GCLOUD_REGION=asia-northeast1
    ```

2. Cloud Spanner の用意

    ```sh
    gcloud beta spanner instances create "${APP_NAME}" \
        --description "Jaguer Lottery" \
        --config "regional-${GCLOUD_REGION}" \
        --processing-units 100
    gcloud spanner databases create app --instance "${APP_NAME}" \
        --ddl='CREATE TABLE swags (id INT64, name STRING(100), stock INT64) PRIMARY KEY(id)'
    ```

3. サービスアカウントの作成

    ```sh
    gcloud iam service-accounts create "${APP_NAME}"
    export SA_EMAIL=$( gcloud iam service-accounts list \
        --filter="name:${APP_NAME}" --format "value(email)")
    gcloud projects add-iam-policy-binding "$( gcloud config get-value project )" \
        --member="serviceAccount:${SA_EMAIL}" --role='roles/spanner.databaseUser'
    gcloud iam service-accounts keys create creds.json \
        --iam-account="${SA_EMAIL}"
    pushd api && ln -s ../jaguer-cn-creds.json creds.json && popd
    pushd web && ln -s ../jaguer-cn-creds.json creds.json && popd
    ```

    書き込む対象の Google Sheets には  
    ここで新規に作成されたメールアドレスに対し編集権限を付与します。

4. API サーバ（Cloud Run）のデプロイ

    ```sh
    gcloud run deploy "${APP_NAME}-api" --region "${GCLOUD_REGION}" \
        --platform managed --allow-unauthenticated --source api \
        --set-env-vars PROJECT_ID="$( gcloud config get-value project )" \
        --set-env-vars SPANNER_INSTANCE="${APP_NAME}" \
        --set-env-vars SPANNER_DATABASE="app" \
        --set-env-vars SPREAD_SHEET_ID="<Google Sheets の ID>" \
        --set-env-vars SPREAD_SHEET_TAB_NAME="<Google Sheets のタブ名>" \
        --service-account "${SA_EMAIL}"
    API_ENDPOINT=$( gcloud run services describe "${APP_NAME}-api" \
        --region "${GCLOUD_REGION}" --format 'value(status.address.url)')
    open "${API_ENDPOINT}/version"
    ```

5. Web サーバ（Cloud Run）のデプロイ

    ```sh
    IMAGE_NAME="gcr.io/$( gcloud config get-value project )/jaguer-cn-lottery:web"
    gcloud builds submit ./web --pack image="${IMAGE_NAME}"
    gcloud run deploy "${APP_NAME}-web" --region "${GCLOUD_REGION}" \
        --platform managed --allow-unauthenticated \
        --image "${IMAGE_NAME}" \
        --set-env-vars API_ENDPOINT="${API_ENDPOINT}"
    SERVICE_URL=$( gcloud run services describe "${APP_NAME}-web" \
        --region "${GCLOUD_REGION}" --format 'value(status.address.url)')
    open "${SERVICE_URL}"
    ```

    できあいのイメージを使う場合

    ```sh
    gcloud run deploy lottery --region asia-northeast1 \
        --image "gcr.io/pottava/jaguer-cn-lottery:web" \
        --set-env-vars API_ENDPOINT="https://jaguer-cn-lottery-api-qbgp4i3oeq-an.a.run.app" \
        --platform managed --allow-unauthenticated
    open "$( gcloud run services describe lottery --region asia-northeast1 \
        --format 'value(status.address.url)')"
    ```
