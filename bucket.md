**upload to cdn**

```bash
rclone copy ./builds/ cdn:cdn/manifest/ --create-empty-src-dirs --progress
```

**delete from cdn**

```bash
aws s3 rm s3://cdn/manifest/<buildid>/ --recursive --endpoint-url https://<accountid>.r2.cloudflarestorage.com
```

**download from cdn**

```bash
chunker.exe download https://cdn.retrac.site/manifest/<manifest> <outputdir>
```

- `++Fortnite+Release-14.40-CL-14550713-Windows.acidmanifest`

**abort multipart uploads**

```bash
for upload in $(aws s3api list-multipart-uploads --bucket cdn --query "Uploads[].{Key:Key, UploadId:UploadId}" --output json --endpoint-url https://<accountid>.r2.cloudflarestorage.com | jq -c '.[]'); do
    key=$(echo $upload | jq -r '.Key')
    upload_id=$(echo $upload | jq -r '.UploadId')
    aws s3api abort-multipart-upload --bucket cdn --key "$key" --upload-id "$upload_id" --endpoint-url https://<accountid>.r2.cloudflarestorage.com
    echo "Aborted upload: $upload_id for $key"
done
```
