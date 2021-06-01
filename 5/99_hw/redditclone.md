Backend для клона реддита.

-----

Я нашел замечательный клон реддита ( https://www.reddit.com/ )

* https://asperitas.vercel.app/ - веб-версия
* https://github.com/d11z/asperitas - исходники ( все на JS ), на всякий случай прикладыю архив с ними в 5/99_hw/code/asperitas-master.zip

К счастью, фронтент дам полностью работает на JS и требует от вас только отдать нужную статику: хтмл, который подключит весь JS, а последний уже будет делать AJAX-запросы к беку.

Бек и предстоит реализовать. Все апи будет отдавать данные в формате JSON.

У вас будут следующие апи-методы:

1) POST /api/register - регистрация
2) POST /api/login - логин
3) GET /api/posts/ - список всех постов
4) POST /api/posts/ - добавление поста - обратите внимание - есть с урлом, а есть с текстом
5) GET /api/posts/{CATEGORY_NAME} - список постов конкретной категории
6) GET /api/post/{POST_ID} - детали поста с комментами
7) POST /api/post/{POST_ID} - добавление коммента
8) DELETE /api/post/{POST_ID}/{COMMENT_ID} - удаление коммента
9) GET /api/post/{POST_ID}/upvote - рейтинг поста вверх - ГЕТ был сделан автором оригинального фронта, я пока не добрался форкнуть и поправить.
10) GET /api/post/{POST_ID}/downvote - рейтинг поста вни
11) GET /api/post/{POST_ID}/unvote - отмена ( удаление ) своего голоса в рейтинге
12) DELETE /api/post/{POST_ID} - удаление поста
13) GET /api/user/{USER_LOGIN} - получение всех постов конкртеного пользователя

Внутри у вас будут следующие сущности:

1) Пользователь
2) Сессия ( получается при авторизации )
3) Пост
4) Коммент к посту

Требования:
1) данные хранятся в памяти, везде есть мютексы
2) полная работоспоосбность всего функционала при запуске приложения
3) корректная работа от разных пользователей
4) глобальные переменные для хранения данных использовать нельзя, все хранится в полях структур
5) Внешние фреймворки (echo, gin и тд) использовать нельзя
6) надо использовать структуру проекта из https://github.com/golang-standards/project-layout

Это задание позволит уже как-то поработать над архитектурой.
Кто хочет идти вперед - можно хранить данные через паттерн репозиторий и не в памяти, а в mysql, а еще сделать сразу же тесты. Это будет тема следующего задания :) Но если вы сделаете кривую архитектуру, а потом улучшите - это наоборот закрепит понимание темы.

Как смотреть в каком формате фронтенд ожидает ответ от апи:

0) https://s.mail.ru/FXHM/33xRndiyK
1) Открыть сайт https://asperitas.vercel.app/
2) Открыть консоль ( в Chrome - F12 )
3) Выбрать секцию Network
4) Выбрать XHR - это покажет только аякс-запросы
5) Кликнуть на запрос
6) В секции Headers->Request Headers при постинге коммента можно найти заголовок authorization - так клиент отправляет авторизацию
7) В секции Preview можно найти ответ

Смотрим ответ, пишем код, который будет отдавать точно такой же ответ.

* фронтенд-часть находится в папке redditclone/template, index.html надо отдать гошным сервисом корне ( / ), js и css - как статику, тоже гошным кодом
* Отдача шаблона и статики были в 4-й лекции
* В качестве роутинга можно использовать gorilla/mux (обратите внимание как рабосттаь со статикой там - смотрите доку в гитлабе)

Тому, как делать это задание посвящен 1-й вебинар (crudapp). Так же эти темы разбираются в 3-й части курса. Как работать с JWT можно посмотреть в 1-й недели 3-й части курса, рекомендую глянуть соответствующие лекции. Это задание даст вам очень много опыта по разработке чего-то более приближенного к реальности и потому оно подвинуто вперед, а не оставлено под самый конец. Чтобы пока есть время - все успели его сделать и прокачаться. В JWT работайте не с мапой, а со структурой на упаковку-распаковку.

Напоминаю, что код решания надо писать САМОСТОЯТЕЛЬНО. Не копипастите из crudapp (вебинар) или photolist (3-я часть). Я хочу чтобы вы научились все это делать сами, а не полагаться на то что было в курсе.

Вендоринг коммитить не надо. Гитлабу плохо от такого количество обновленных файлов и ревью превращается в ад. Но это касается конкретно домашек с курса, а не прода.