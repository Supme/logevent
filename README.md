# Logevent
Windows log event metric for Prometheus textfile inputs

Для логирования неуспешных входов:
 - включить логирование плохих входов в windows через gpedit.msc
   Local Computer Policy → Computer Configuration → Windows Settings → Security Settings → Local Policies → Audit Policy
   
   или по русски
   
   Политика "Локальный компьютер" → Конфигурация компьютера → Конфигурация Windows → Параметры безопасности → Локальные политики → Политика аудита
   поставить обе галки в Audit logon events (Аудит входа в систему) и Audit account logon events (Аудит событий входа в систему)
   
 - Скопировать logevent.exe в C:\Program Files\windows_exporter\logevent.exe

 - В планировщике импортируем Security.xml (если это сделали, то следующий шаг пропускаем)
 
 - В планировщике заданий создаём задачу в "Задачи просмотра событий":
   Имя "Security"
   Выполнять с наивысшими правами
   Выполнять вне зависимости от регистрации пользователя
   Триггер создать:
   - начать задачу: при событии
   - при событии настраиваемое
   - фильтр событий по журналу Журнал событий Журналы windows безопасность
   - коды событий: 4625,5461,529,530,531,532,533,534,535,539
   Действия:
     Запуск программы:
	 "C:\Program Files\windows_exporter\logevent.exe" аргументы -m "windows_logevent_bad_login_count" -d "Bad login event"
   Параметры:
   - Останавливать задачу выполняемую дольше: 1ч.
   - Если задача уже выполняется, то применять правило: Запускать новый экземпляр задания
   
У windows-exporter должен быть включен коллектор textfile

Запрос в Prometheus примерно таков:
increase(windows_logevent_bad_login_count[5m])
