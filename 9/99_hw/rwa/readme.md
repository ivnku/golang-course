Тема: API для приложения RealWorld

Проблема придумать хорошее учебное приложение на самом деле довольно актуальна :)

Почти всегда это в том или ином виде CRUD. Но просто CRUD над одной сущностью это не так интересно. Поэтому надо чтобы сущностей было побольше со всякими связими между ними. Это позволит получить как раз опыт архитектуры приложения и подхода с паттерном репозиторий.

Есть на просторах интернета репозиторий https://github.com/gothinkster/realworld . В нем пишется клон проекта Medium на различных языках, как фронтентд, так и бакенд. Называется Conduit. Посмотрель в реальности его можно на сайте https://demo.realworld.io/#/ .

Как вы догадываетесь, мы будем писать свой бакенд для этого приложения.

Вернее, какое-то его подмножество - я сделал в тесте немного меньше сущностей, чем есть в самом realworld. +2 сущности радикально новому ничему не научат, а время потратят.

У самого realworld есть набор тестов через Postman ( https://github.com/gothinkster/realworld/tree/master/api ) - но меня не удовлетворили те тесты, поэтому я реализовал части из них на гошке, чтобы вам не пришлось ставить никакую ноду.

В рамках этого задания вам надо попробовать:
1. Разбиение проекта на отдельные компоненты (хендлеры-репозитории)
2. Сессии (передаются через хедер Authorization) должны быть stateful. Если вы решите использовать JWT - то все равно сессия должна быть stateful. Причем храниться должна нормально, а не так что запомнили весь токен.

Реализовать придется следующие сущности:
* Юзер
* Сессия
* Статьи с различными фильтрами

По желанию, конечно же можнео реализовать все остальные сущности с отдельными тестами.

Можно следовать предлагаемой структуре и не бить по пакетам, а можно чуть подправить и разделить как если бы это было бы в реальности - см crudapp за основу. Код надо писать в файле realworld.go если вы решите следовать 1-му варианту.

Как обычно во всех заданиях этого и предыдущих курсов - данное описание - это базовая постановка задачи, все остальное придется получать из тестов. Тестов не очень много, но они выполнены достаточно универстально в табличном виде и некоторым количеством магии, разобраться в которой вам так же будет полезно :)

Так же у проекта realworld есть swagger схема. Она приложена в задании в папке swagger. Можно запустить сервер документации (находясь в папке rwa/swagger ):

* go get -u github.com/go-swagger/go-swagger/cmd/swagger
* swagger serve swagger.json -p 8085

Успехов!

P.S. В этой домашке могут быть баги, будьте внимательны. Но у меня есть рабочее решение, так что убедитесь сначала что это именно бага, а вопросы по реализации.

Как работают тесты:

* В этом задании у вас есть большой набор интеграционных тестов. Интеграционных это значит что они тестируют всю цепочку с учетом изменяемого состояния системы - если я добавил что-то, то потом должен уметь это прочитать. 
* За счет этого вы можете двигаться буквально по одному шагу, дописывая код чтобы выполнился очередной тест-кейс - как это было во 2-й домашке с игрой.
* Все тесты располагаются в app_test.go - там от вас получается хттп-хендлер (см план работы ниже) и запихивается в тестовый сервер.
* у теста есть имя - по нему обычно понятно что мы тестируем. Остальыне поля +- говорят за себя. Обратите внимание - кое-где есть триггеры "до" и "после" - в них происходит, например, установка-замена авторизационного токена.

Про токен:

Токен это, в общем, сессия. Но которая передается не через куки, а черех хттп-хедеры. Вы можете или использовать токен как сессионный ключ, или сделать jwt-токен с доп дарными (можете в redditclone про это прочитать). Рекомендую начать с сессионного ключа (который у вас будет просто ключем в мапке сессий).

Про POST&ko запросы:

На сервер уходит тело запроса в формате жсон, не форм-урлэнкодед. Значит надо вычитать тело и распаковать жсон. Обязательно проверять ошибки! В предыдущей лекции было как вычитывать тело запроса.

План работы:

* С точки зрения кода - вся 4-я лекция посвящена тому как делать это задание. но не копипастите код оттуда!!! все писать самостоятельно 
* В результате вам надо раздробить решение на пакеты в соотвтетсвии с https://github.com/golang-standards/project-layout и тем что я рассказывал на лекции с crudapp
* Тут нет параметров в урле, но есть разделение по GET/POST/etc методам - можно или через switch-case зайти в нужную функцию, или вкрутить другой роутер (например, gorilla/mux или что побыстрее) и сразу роут цеплять к нужному методу. в лекции был пример. но обратите внимание как в gorilla/mux использовать миддлверы - почитайте там в доке
* Помните что нам надо разбить логику на слои: handler -> repository -> db (сейчас слайсы и мапы, в след части домашки - реальные базы).