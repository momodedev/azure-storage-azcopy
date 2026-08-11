from azure.storage.blob import BlobServiceClient, BlobClient, ContainerClient 
import oss2
import logging

logging.basicConfig(level=logging.DEBUG)

# Aliyun OSS SDK for Python 
def migrate_files(aliyun_access_key_id, aliyun_access_key_secret, aliyun_endpoint, aliyun_bucket_name, azure_connection_string, azure_container_name): 
# Initialize Aliyun OSS 
    auth = oss2.Auth(aliyun_access_key_id, aliyun_access_key_secret) 
    bucket = oss2.Bucket(auth, aliyun_endpoint, aliyun_bucket_name) 
# Initialize Azure Blob Storage 
    blob_service_client = BlobServiceClient.from_connection_string(azure_connection_string) 
    container_client = blob_service_client.get_container_client(azure_container_name) 
# List objects in Aliyun OSS bucket 
    for obj in oss2.ObjectIterator(bucket): 
        if obj.is_prefix(): 
            continue 
    # skip directories 
        else: object_name = obj.key 
        aliyun_object_url = f"https://{aliyun_bucket_name}.{aliyun_endpoint}/{object_name}" 
# Make sure you generate a valid SAS URL or use any kind of authorization method needed 
# Example here expects you have direct URL access 
# Read from Aliyun OSS (We only need the URL to be used by Azure) 
# Azure's Put Block From URL API makes an HTTP(S) GET request to fetch the content from Aliyun OSS 
        blob_client = container_client.get_blob_client(blob=object_name) 
        print(f"Transferring {object_name} from Aliyun OSS to Azure Blob Storage") 
        blob_client.stage_block_from_url(block_id='1', source_url=aliyun_object_url) 
# Commit the Block (property '1' used in stage_block_from_url should match) 
        blob_client.commit_block_list(['1']) 
        print("Migration completed successfully.") 
# Aliyun OSS credentials 
ALIYUN_ACCESS_KEY_ID = "LTAI5tHi3ARywbp1rYehbyRz" 
ALIYUN_ACCESS_KEY_SECRET = "i0N3Ek2BbET7euG9c8y7uGkwyg0504" 
ALIYUN_ENDPOINT = "oss-eu-central-1.aliyuncs.com" 
ALIYUN_BUCKET_NAME = "dmo0001" 
# Azure Blob Storage credentials 
AZURE_CONNECTION_STRING = "DefaultEndpointsProtocol=https;AccountName=s32blobdestination01;AccountKey=7NumulCKotPFGR7qq3Ung12P4cKR845/yhxZ1CHO5syOLaXIgJjRp8jfoViDdsphqdHE2yPyvZd9+ASt+8Iirg==;EndpointSuffix=core.windows.net" 
AZURE_CONTAINER_NAME = "20241030" 
# Running the migration 
migrate_files(ALIYUN_ACCESS_KEY_ID, ALIYUN_ACCESS_KEY_SECRET, ALIYUN_ENDPOINT, ALIYUN_BUCKET_NAME, AZURE_CONNECTION_STRING, AZURE_CONTAINER_NAME) 