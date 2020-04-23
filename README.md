# go-large-file-encrypt

该程序可以使用RSA算法对超大文件进行非对称加密，并将加密后的文件打包上传至aws s3; 使用yubikey的私钥一键将加密文件解密。
测试时加密一个10G的测试文件需要40m，aws s3分片上传需要40m。
特点：
* 通过混合加密提供RSA加密的安全性和AES加密的性能，10G文件加密<1h
* AES加密不用加载整个待加密文件，内存占用可控
* 支持多个公钥GPG公钥加密AES秘钥
* aws s3分片上传，支持最大5T单个文件

加密流程：
* 生成随机AES密钥
* 通过GPG RSA加密算法加密AES秘钥，AES算法加密待加密文件
* 计算原文件的sha256
* AES加密大文件,使用CFB算法,单线程
* ZIP打包加密后的备份文件、加密AES秘钥、sha256
* 上传打包文件至s3 
* 清理临时文件

解密流程：
* GPG解密出aes key
* AES解密出原文件

this is a program for encrypt/decrypt large files, encrypt/decrypt original file using AES and encrypt/decrypt AES secrets using RSA. Suitable for very large files, because not reading them to the memory.

