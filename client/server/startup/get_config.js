/**
 * Created by Than on 31.07.2016.
 */

// Получаем календарь
getData(getCalendar).then(
    periods => {
        calendar = periods;
        // Пустой промис, чтобы гарантированно дождаться выполнения текущего шага + нотификация об успешности
        return new Promise((resolve) => {
            console.log("Календарь успешно загружен");
            return resolve();
        });
    }
).catch(
    error => {
        console.log(error);
    }
);