import requests 
from azure.storage.blob import BlobServiceClient, BlobClient 
from oss2 import Auth, Bucket 

# Aliyun OSS credentials a
liyun_access_key_id = 'your-access-key-id' 
aliyun_access_key_secret = 'your-access-key-secret' 
aliyun_endpoint = 'oss-eu-central-1.aliyuncs.com' 
# Modify based on your bucket's location 
aliyun_bucket_name = 'dmo0001' 
# Azure storage account details 
connect_str = "" 
container_name = "20241030" 
# Create an instance of the Aliyun Auth class 
auth = Auth(aliyun_access_key_id, aliyun_access_key_secret) 
# Create an instance of the Aliyun Bucket class 
bucket = Bucket(auth, aliyun_endpoint, aliyun_bucket_name) 
# Specific file details 
object_name = 'icudtl.dat' 
azure_blob_name = 'icudtl.dat' 
# Generate a signed URL for the Aliyun OSS object 
aliyun_file_url = bucket.sign_url('GET', object_name, 60) 
# URL valid for 60 seconds 
# Instantiate a BlobServiceClient using a connection string 
blob_service_client = BlobServiceClient.from_connection_string(connect_str) 
container_client = blob_service_client.get_container_client(container_name) 
blob_client = container_client.get_blob_client(azure_blob_name) 
# Download content from the Aliyun OSS URL 
response = requests.get(aliyun_file_url) 
if response.status_code == 200: 
# Upload content to Azure Blob Storage 
    blob_client.upload_blob(response.content, overwrite=True) 
    print("File migrated successfully!") 
else: print("Failed to download file from Aliyun OSS.")