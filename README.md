Про репозиторий:

1. Репозиторий должен быть приватным. Пожалуйста, не выкладывайте ваши решения в публичный доступ ( только на степик в соответствующий раздел обмена решениями ). Я не хочу отчислять учащихся за то что они сдают ваш код, прецеденты были.
2. Форкаем в гитлабе, выдаем доступ пользователю *rvasily* в ролью Maintainer. Не romanov.vasily ( это тоже мой, но другой ). У неправильного акка специально стоит аватарка со знаком STOP.


Про домашки:

0. Репозиторий - первичный источник правды про домашки. Смотрите в первую очередь сюда, а не на степик.
1. Домашки сложные. Халявы не будет, придется попотеть. Но зато вы потом реально будете знать язык.
2. Некоторые вещи могут показаться вам нелогичными, но это сделано специально, чтобы пописать побольше кода. Например некоторые команды в игре ( ДЗ-1 ).
3. Домашки можно пропускать если не нравится.
4. Не надо искать решения домашек и тем более сдавать списанные домашки (увы, были прецеденты). Вы тут чтобы получить знания и научиться решать эти задачи самостоятельно (ну, с моей помощью), а не гуглить их. Если что-то непонятно - надо спрашивать меня, а не пытаться обойти
5. Копипастить код из репы курса не надо. По возможности старайтесь все реализовывать с 0. Так вы научитесь решать эти задачи, а не копи-пастить чужое. Лучше потом вы скопипастите свое, но вы будете понимать почему оно было сделано именно так и как оно работает.
6. Вендоринг коммитить не надо. Гитлабу плохо от такого количество обновленных файлов и ревью превращается в ад. Но это касается конкретно домашек с курса, а не прода.
7. НАДО задавать вопросы. Не тупите с аргументацией "неудобно было вас отвлекать". Я тут специально чтобы вы меня отвлекали вопросами. Если мне неудобно отвечать сразу - я отвечу как смогу. Лучше вы потратите 5 минут моего времени, чем 3 дня будете сидеть сами. В разумных пределах конечно же, задачу я за вас решать не буду, но направление всегда подскажу.

Порядок выполнения домашки:

1. забираем обновление репозитория
2. пишем код
3. тестируем локально (в каждой домашке есть список всех команд)
4. коммитим в отдельную ветку
5. создаем МР из этой ветки в мастер. В МР гораздо удобнее комментировать и видны все изменения ваши
6. на каждый коммит запускается gitlab-ci, он проверяет что ваш код компилится, код отформатирован, тесты работают, все хорошо. в интерфейсе видно вот тут: https://s.mail.ru/LsUU/YaCvsy7h1 + придёт письмо на почту.
7. как загорелось зеленым - можно писать преподавателю с просьбой поревьювить
8. пройденные тесты - это приглашение обсудить вашу работу на код ревью, а не окончание работы.
9. грузим решение на степик чтобы засчиталась оценка, так же можно загрузить все в раздел с решениями чтобы увидеть решения других участников
10. Merge Request в основную ( мою ) репу создавать не надо, работайте только в своей приватной репе
11. Если будут правки - не надо каждый фикс оформлять в виде отдельного коммита так что их будет 10 штук на ревью. Это тяжело ревьбвить потом - я смотрю что вы коммитили.

Для тех кто хочет дополнительно прокачаться - записывайте после каждой недели видео на телефон, о том чему вы научились, что как работает. Как если бы вы отвечали про это на собеседовании. Это структурирует полученную в голове информацию, чтобы вы сам могли объяснить ее. А видео присылайте мне. Можно, например, через cloud.mail.ru. Не все буду сомтреть в полном объеме, но мельком прогляжу.

Как забирать обновления из основной репы:
```
git pull https://gitlab.com/rvasily/golang-stepik-2021q2.git
git push origin master
```