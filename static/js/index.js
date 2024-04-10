// Деавторизация
const unauth = () => {
    console.log('unauth')
    axios.post('/unauth', {}).then((response) => {
        console.log('Response: ' + response)
        window.location.reload()
    }).catch((error) => {
        console.log('Error: ' + error)
    })
}

// Загрузка страницы
const pagination = () => {
    axios.post('/api/get_page', {
        offset: 0,
        limit: 10
    }).then((response) => {
        console.log(response)
    }).catch((error) => {
        console.log(error)
    })
}

// Отправка файла
const uploadFile = () => {
    // Получаем форму по ID
    const form = document.getElementById('uploadForm')

    const file = fileInput.files[0]
    if (!file) {
        // Если файл не выбран, выводим сообщение об ошибке и завершаем функцию
        const text = 'Файл не выбран'
        alert(text)
        console.error(text)
        return
    }

    // Создаем объект FormData и добавляем файл из формы
    const formData = new FormData(form)

    // Отправляем файл на сервер с помощью Axios
    axios.post('/api/upload', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    }).then((response) => {
        // Обрабатываем успешный ответ, если необходимо
        console.log('Файл успешно загружен', response)
        form.reset()
    }).catch((error) => {
        // Обрабатываем ошибку, если необходимо
        const text = 'Ошибка при загрузке файла'
        alert('Ошибка при загрузке файла', error)
        console.error('Ошибка при загрузке файла', error)
    })
}

const columns = [
    { data: 'Томская область' },
    { data: 'Регион' },
    { data: 'Назначен ответственный' },
    { data: 'Страница подтверждена' },
    { data: 'Ссылка на официальную страницу Вконтакте' },
    { data: 'Ссылка на официальную страницу Одноклассники' },
    { data: 'Ссылка на официальную страницу Telegram' },
    { data: 'Официальная страница не ведется на основании' },
    { data: 'Комментарий по НПА' },
    { data: 'Полное наименование объекта' },
    { data: 'ОГРН' },
    { data: 'Статус' },
    { data: 'Комментарий' },
]

$(document).ready(function () {
    $('#mainTable').DataTable(
        {
            // data: dataSet,
            columns: columns,
            ajax: '/api/table',
            serverSide: true,
        }
    );
});
