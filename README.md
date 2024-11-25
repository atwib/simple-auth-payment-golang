## **Payment Application Backend - Technical Test**

Backend Developer Role | Agung Tri Wibowo  
Project Duration: 22 November 2024 - 26 November 2024  
Tech Stack:

-   Golang
-   Gin Framework
-   JSON

----------

Application Overview:  
Aplikasi ini adalah layanan pembayaran yang mendukung transfer antara customer dan merchant dengan skenario berikut:

-   Customer ke Customer
-   Customer ke Merchant
-   Merchant ke Merchant
-   Merchant ke Customer

Features and User Flows:

1.  User Registration and Login:
    
    -   Pengguna mendaftar untuk membuat akun.
    -   Pengguna login untuk otentikasi menggunakan kredensial.
2.  Wallet Registration:
    
    -   Pengguna mendaftar sebagai Customer atau Merchant, yang akan memengaruhi aturan transaksi (biaya, saldo minimum, dll.).
3.  Top-Up Wallet:
    
    -   Pengguna dapat menambahkan dana ke wallet mereka.
4.  Transfer Antara Wallet:
    
    -   Mendukung transfer aman antar-wallet berdasarkan jenis pengguna dan biaya yang berlaku.
5.  Konfigurasi/Aturan:
    
    -   Minimal Top-Up: 10,000
    -   Saldo Minimum Customer: 5,000
    -   Saldo Minimum Merchant: 0
    -   Biaya Admin Transfer: 1,000
    -   Biaya Withdraw Merchant: 0
    -   Biaya Withdraw Customer: 1,000

----------

API Endpoints:

Authentication Endpoints:

1.  Register User
    
    -   Endpoint: POST [`http://localhost:8080/api/auth/register`](http://localhost:8080/api/auth/register)
    -   Request Body:  
    {  
    "username": "username",  
    "password": "password"  
    }

2.  Login User
    
    -   Endpoint: POST [`http://localhost:8080/api/auth/login`](http://localhost:8080/api/auth/login)
    -   Request Body:  
    {  
    "username": "username",  
    "password": "password"  
    }

3.  Get Current User
    
    -   Endpoint: GET [`http://localhost:8080/api/auth/me`](http://localhost:8080/api/auth/me)
    -   Header: Authorization: Bearer
4.  Refresh Token
    
    -   Endpoint: POST [`http://localhost:8080/api/auth/refresh-token`](http://localhost:8080/api/auth/refresh-token)
-   Request Body:  
    {  
    "refresh_token": "<refresh_token>"  
    }

5.  Reset Password
    
    -   Endpoint: POST [`http://localhost:8080/api/auth/reset-password`](http://localhost:8080/api/auth/reset-password)
    -   Request Body:  
    {  
    "old_password": "current_password",  
    "new_password": "new_password"  
    }

6.  Logout User
    
    -   Endpoint: POST [`http://localhost:8080/api/auth/logout`](http://localhost:8080/api/auth/logout)
    -   Header: Authorization: Bearer

Wallet Endpoints:

1.  Register Wallet
    
    -   Endpoint: POST [`http://localhost:8080/api/wallet/register`](http://localhost:8080/api/wallet/register)
    -   Request Body:  
    {  
    "type": "Merchant",  
    "pin": 1234,  
    "pin_confirm": 1234  
    }

2.  Top-Up Wallet
    
    -   Endpoint: POST [`http://localhost:8080/api/wallet/top-up`](http://localhost:8080/api/wallet/top-up)
    -   Request Body:  
    {  
    "amount": 50000  
    }

3.  Transfer Confirmation
    
    -   Endpoint: POST [`http://localhost:8080/api/wallet/transfer-confirm`](http://localhost:8080/api/wallet/transfer-confirm)
    -   Request Body:  
 

   {  
    "to_wallet_id": "450f0611-dbb9-4313-9443-e351d00645d4",  
    "amount": 10000  
    }

4.  Transfer Funds
    
    -   Endpoint: POST [`http://localhost:8080/api/wallet/transfer`](http://localhost:8080/api/wallet/transfer)
    -   Request Body:  
{  
"to_wallet_id": "450f0611-dbb9-4313-9443-e351d00645d4",  
"amount": 10000,  
"pin": 1234  
}

----------

Project Setup:

1.  Clone the Repository:  
    git clone  https://github.com/atwib/simple-auth-payment-golang.git
    
2.  Install Dependencies:  
    go mod tidy
    
3.  Run the Application:  
    go run main.go
    
4.  Access the API:  
    The application will be accessible at [http://localhost:8080](http://localhost:8080/)
    

----------

Notes:

-   Pastikan lingkungan Go sudah dikonfigurasi dengan benar.
-   Gunakan REST client seperti Postman atau cURL untuk pengujian.
-   Authentication memerlukan bearer token setelah login.

Developed by Agung Tri Wibowo  
For MNC Technical Test