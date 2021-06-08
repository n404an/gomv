# gomv

gomv - это надстройка над системной утилитой mv, написанная на Golang

## Возможности:
* Множество исходных папок
* Множество целевых папок
* Маска для исходных файлов
* Кол-во параллельных операций
* Скан исходных папок раз в минуту в поисках новых файлов.

## Установка

[Тут](https://github.com/n404an/gomv/releases/latest)
готовый бинарник gomv
(права на выполнение нужно добавить)
```bash
chmod +x gomv
```

или из исходников
```bash
go build cmd/main.go
```

## Использование
На примере:
```bash
./gomv -src test/src/[1-4] /mnt/hdd1 /home/sdc1/tmp -dst test/dst/* /mnt/targetHdd -p 4 -m 1[5-8]*.log
```
Будут сканироваться папки test/src/1 test/src/2 test/src/3 test/src/4 /mnt/hdd1 /home/sdc1/tmp

по маске 1[5-8]*.log (по-умолчанию *)

всё, что проходит по критериям встаёт в очередь на перенос в папки назначения test/dst/*(все папки в каталоге dst) /mnt/targetHdd

в 4 потока  (по-умолчанию 1)

## Будущие возможности:
* Поддержка Windows.
* Определение доступного объёма в папке назначения, если забита - пропустить папку.
* Определение физического носителя, что бы исключить параллельное копирование с него или на него.
* Определение сетевого диска, что бы учесть нерезиновый канал.
* Последовательный режим обхода папок назначения.
* Прибраться в коде.

## Задонатить
Если тебе зашла прога и ты хочешь, что бы она дальше развивалась, можешь меня немножко поощрить какой-нибудь монеткой. :)

BTC [bc1qze70cnewu3rx2msnv7vnpkf34dzj409lzdcgtl](https://www.blockchain.com/btc/address/bc1qze70cnewu3rx2msnv7vnpkf34dzj409lzdcgtl)

ETH или любой ERC20 токен [0x913c6dfeB21eB91121a9190fA3661ba60Ce06b81](https://etherscan.io/address/0x913c6dfeb21eb91121a9190fa3661ba60ce06b81)

XCH [xch1kjdysu80782g234fh7skptlm30h2q98907nd5674vyaseacrujxqz24rdz](https://www.chiaexplorer.com/blockchain/address/xch1kjdysu80782g234fh7skptlm30h2q98907nd5674vyaseacrujxqz24rdz)
