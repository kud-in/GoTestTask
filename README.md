# Test Task for Golang Developer

Требуется разработать систему, которая будет получать, хранить и отдавать обработанные данные. Данные, которые необходимо получать — это курсы валют (USD, EUR) из открытого источника на ваш выбор. Необходимо получать и сохранять данные несколько раз в минуту в базе данных (на ваш выбор). Далее данные могут быть получены из данного сервиса по API.
Необходимые запросы:

1) Получение статуса, возвращает текущее значение валюты (последнее), так же должен
отдать среднюю стоимость за 24 часа, за неделю, за месяц в одном запросе.

2) Получение истории агрегированных данных. Запрос возможен с ограничением периода запрашиваемых данных, всегда запрашивается вид агрегации (среднее за 1 минуту, за 5
минут, за 1 час, за день). Ответ запроса содержит массив элементов являющимися единицами одного вида агрегации, среднюю стоимость валюты за каждый момент времени (08-09-2019 17:00 – 65.4, 08-09-2019 17:05 – 64.3).

3) Запрос стоимость валюты за момент времени. Параметром запроса является отметка времени, результатом значение стоимости.

* Приложение можно дробить на необходимые компоненты.
* Структуру кода и уровень детализации исполнения вы выбираете сами.
* Базу данных и инструменты для работы вы выбираете на своё усмотрение.
* Форматы ответов не регламентированы, строите исходя из своего видения.
* Результат работы требуется залить в репозиторий на GitHub.
* Требуется отследить и указать потраченное время на реализацию задачи.

### Run App

```
git clone https://github.com/kudin-0x57/GoTestTask.git
cd ./GoTestTask
docker network create -d bridge GoTestNetwork
docker-compose -f docker-compose.yml up --build -d
```

### Сlick on the link [http://localhost:8080](http://localhost:8080)