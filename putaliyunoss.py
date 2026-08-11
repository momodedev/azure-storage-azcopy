from azure.storage.blob import BlobServiceClient, BlobClient

# Replace the following with your Azure storage account details 
connect_str = "DefaultEndpointsProtocol=https;AccountName=s32blobdestination01;AccountKey=7NumulCKotPFGR7qq3Ung12P4cKR845/yhxZ1CHO5syOLaXIgJjRp8jfoViDdsphqdHE2yPyvZd9+ASt+8Iirg==;EndpointSuffix=core.windows.net" 
container_name = "20241030" 
# Instantiate a BlobServiceClient using a connection string 
blob_service_client = BlobServiceClient.from_connection_string(connect_str) 
container_client = blob_service_client.get_container_client(container_name) 
# Aliyun OSS file URL and the Azure Blob name where it will be stored 
aliyun_file_url = "https://dmo0001.oss-eu-central-1.aliyuncs.com/icudtl.dat" 
azure_blob_name = "icudtl.dat" 
# Instantiate a BlobClient 
blob_client = container_client.get_blob_client(azure_blob_name) 
# Upload content from the Aliyun OSS URL 
blob_client.upload_blob_from_url(aliyun_file_url, overwrite=True) 
print("File migrated successfully!") 