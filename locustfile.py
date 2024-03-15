from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(5, 15)
    
    def on_start(self):
        tokenResponse = self.client.post("/v1/user/login", json={
            "username": "usertest1",
            "password": "testpassword"
        })
        response = tokenResponse.json()

        responseToken = response['data']['accessToken']
        self.client.headers = {'Authorization': 'Bearer ' + responseToken}

        bankAccount = self.client.post("/v1/bank/account", json={
            "bankName": "testBankName",
            "bankAccountName": "Test bank name",
            "bankAccountNumber": "adsjasdkjasldjaskljd123"
        })
        response = bankAccount.json()
        self.bankAccountId = response['data']['id']

        product = self.client.post("/v1/product", json={
            "name": "test product",
            "price": 100000,
            "imageUrl": "https://i0.wp.com/dowse.co.uk/wp-content/uploads/2020/03/white-square.png",
            "stock": 10,
            "condition": "second",
            "tags": ["test"],
            "isPurchasable": False
        }) 
        response = product.json()
        self.productId = response['data']['id']


    
    @task
    def indexProduct(self):
        self.client.get("/v1/product?limit=5&offset=0&tags=test&tags=wew&condition=second&search=test&maxPrice=100000&minPrice=0")
        
    @task
    def detailProduct(self):
        self.client.get("/v1/product/"+self.productId)

    @task
    def createProduct(self):
        self.client.post("/v1/product", json={
            "name": "test product",
            "price": 100000,
            "imageUrl": "https://i0.wp.com/dowse.co.uk/wp-content/uploads/2020/03/white-square.png",
            "stock": 10,
            "condition": "second",
            "tags": ["test"],
            "isPurchasable": False
        })

    @task
    def patchProduct(self):
        self.client.patch("/v1/product/" + self.productId, json={
            "name": "test product",
            "price": 100000,
            "imageUrl": "https://i0.wp.com/dowse.co.uk/wp-content/uploads/2020/03/white-square.png",
            "stock": 10,
            "condition": "second",
            "tags": ["test"],
            "isPurchasable": False
        })


    @task
    def indexBankAccount(self):
        self.client.get("/v1/bank/account")

    @task
    def createBankAccount(self):
        self.client.post("/v1/bank/account", json={
            "bankName": "testBankName",
            "bankAccountName": "Test bank name",
            "bankAccountNumber": "adsjasdkjasldjaskljd123"
        })

    @task
    def patchBankAccount(self):
        self.client.patch("/v1/bank/account/" + self.bankAccountId, json={
            "bankName": "testBankName",
            "bankAccountName": "Test bank name",
            "bankAccountNumber": "adsjasdkjasldjaskljd123"
        })

    @task
    def buyProduct(self):
        self.client.post("/v1/product/" + self.productId + "/buy", json={
            "bankAccountId" : self.bankAccountId,
            "paymentProofImageUrl": "https://i0.wp.com/dowse.co.uk/wp-content/uploads/2020/03/white-square.png",
            "quantity" : 10
        })