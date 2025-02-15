from locust import HttpUser, task, between, SequentialTaskSet
import random
import threading

existing_users = set()
lock = threading.Lock()

class UserBehavior(SequentialTaskSet):
    def on_start(self):
        """
        При старте пользователь проходит аутентификацию через POST /api/auth.
        При первой аутентификации пользователь создаётся автоматически.
        Имя пользователя сохраняется и добавляется в глобальный список для дальнейшего использования в качестве получателя.
        """
        while True:
            self.username = f"user_{random.randint(1, 1000000)}"
            if self.username not in existing_users:
                break

        password = "secret"
        with self.client.post(
            "/api/auth",
            json={"username": self.username, "password": password},
            headers={"Content-Type": "application/json"},
            name="/api/auth",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                self.token = response.json().get("token")
                with lock:
                    existing_users.add(self.username)
            else:
                self.token = None
                response.failure(f"Auth failed: {response.text}")

    @task
    def buy_merch(self):
        """
        Сценарий покупки мерча через GET /api/buy/{item}.
        Перед покупкой запрашивается баланс через GET /api/info.
        Если средств достаточно для покупки хотя бы одного товара, выбирается случайный товар из доступных.
        """
        if not self.token:
            return
        headers = {"Authorization": f"Bearer {self.token}"}

        with self.client.get(
            "/api/info", 
            headers=headers, 
            name="/api/info", 
            catch_response=True
        ) as info_response:
            if info_response.status_code != 200:
                info_response.failure(f"Failed to get info: {info_response.text}")
                return
            info = info_response.json()
            coins = info.get("coins", 0)

        price_map = {
            "t-shirt": 80,
            "cup": 20,
            "book": 50,
            "pen": 10,
            "powerbank": 200,
            "hoody": 300,
            "umbrella": 200,
            "socks": 10,
            "wallet": 50,
            "pink-hoody": 500
        }
        affordable_items = [item for item, price in price_map.items() if coins >= price]
        if not affordable_items:
            return

        item = random.choice(affordable_items)
        with self.client.get(
            f"/api/buy/{item}", 
            headers=headers, 
            name=f"/api/buy/{item}",
            catch_response=True
        ) as response:
            if response.status_code != 200:
                response.failure(f"Buy failed for item '{item}': {response.text}")

    @task
    def send_coin(self):
        """
        Сценарий отправки монет через POST /api/sendCoin.
        Перед отправкой запрашивается баланс через GET /api/info.
        Если баланс достаточен, выбирается случайный получатель из глобального списка (не самого себя)
        """
        if not self.token:
            return
        headers = {"Authorization": f"Bearer {self.token}", "Content-Type": "application/json"}
        
        with self.client.get(
            "/api/info", 
            headers=headers, 
            name="/api/info",
            catch_response=True
        ) as info_response:
            if info_response.status_code != 200:
                info_response.failure(f"Failed to get info: {info_response.text}")
                return
            info = info_response.json()
            coins = info.get("coins", 0)
        
        send_amount = 10
        if coins < send_amount:
            return

        with lock:
            possible_recipients = [user for user in existing_users if user != self.username]
        if not possible_recipients:
            return
        recipient = random.choice(possible_recipients)

        payload = {"toUser": recipient, "amount": send_amount}
        with self.client.post(
            "/api/sendCoin", 
            json=payload, 
            headers=headers,
            name="/api/sendCoin",
            catch_response=True
        ) as response:
            if response.status_code != 200:
                response.failure(f"Send coin failed: {response.text}")

class AvitoShopUser(HttpUser):
    tasks = [UserBehavior]
    wait_time = between(1, 3)
    host = "http://localhost:8080"
