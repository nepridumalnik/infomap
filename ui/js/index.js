// Деавторизация
const unauth = () => {
    console.log("unauth")
    axios.post("/unauth", {}).then((response) => {
        console.log("Response: " + response)
        window.location.reload()
    }).catch((error) => {
        console.log("Error: " + error)
    })
}

// Отправка файла
const uploadFile = () => {
    // Получаем форму по ID
    const form = document.getElementById('uploadForm')

    // Создаем объект FormData и добавляем файл из формы
    const formData = new FormData(form)

    // Отправляем файл на сервер с помощью Axios
    axios.post('/upload', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    }).then((response) => {
        // Обрабатываем успешный ответ, если необходимо
        console.log('Файл успешно загружен', response)
        form.reset()
    }).catch((error) => {
        // Обрабатываем ошибку, если необходимо
        console.error('Ошибка при загрузке файла', error)
    })
}
