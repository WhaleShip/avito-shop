## Запуск

### 1. Создать venv
```sh
python3 -m venv venv
``` 
### 2. Активировать его
```sh
source venv/bin/activate // linux

.\venv\Scripts\activate   // windows
``` 

### 3. Установить locust
```sh
pip install locust
```

### 4. запустить 
через консоль

```sh
// оставил параметры с которыми проводил тест
locust -f stress_test.py --headless --users 1000 --spawn-rate 20 --run-time 3m --csv=report
```
или web интерфейс (доступен на http://0.0.0.0:8089)
```sh
locust -f stress_test.py --headless --users 1000 --spawn-rate 20 --run-time 3m --csv=report
```