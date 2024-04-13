// Переводы
// https://datatables.net/reference/option/language
const language = {
    search: 'Поиск',
    processing: 'Обработка запроса',
    info: 'Показать от _START_ до _END_ из _TOTAL_ записей',
    thousands: '.',
    infoFiltered: '[отфильтровано _MAX_ записей]',
    lengthMenu: '_MENU_ записей на странице',
    zeroRecords: 'Совпадений не обнаружено',
    emptyTable: 'В таблице пока нет записей',
    infoEmpty: 'Записей нет',
    paginate: {
        'first': 'Первая',
        'last': 'Последняя',
        'next': 'Следующая',
        'previous': 'Предыдущая',
    },
    aria: {
        'orderable': 'Сортировать по этому столбцу',
        'orderableReverse': 'Сортировать по этому столбцу в обратном порядке',
    },
}

// Завершение инициализации
const initComplete = () => {
    console.log('initComplete()')
}

// Завершение инициализации
const infoCallback = (settings, start, end, max, total, pre) => {
    // console.log('infoCallback()')
    // console.log('settings: ' + settings)
    // console.log('start: ' + start)
    // console.log('end: ' + end)
    // console.log('max: ' + max)
    // console.log('total: ' + total)
    // console.log('pre: ' + pre)

    const tbody = document.querySelector('#mainTable tbody')
    const rows = tbody.querySelectorAll('tr')

    rows.forEach(function (row) {
        var cells = row.querySelectorAll('td')
        cells.forEach(function (cell) {
            // console.log(cell.textContent)
        })
    })
}

// Функция для удаления строки
const deleteRow = () => {
    // Реализация удаления строки
    console.log('Удаление строки')
}

// Функция для добавления строки
const addRow = () => {
    // Реализация добавления строки
    console.log('Добавление строки')
}

// Обработка загрузки страницы
const init = () => {
    getTable().then((response) => {
        const ins = $('#mainTable').DataTable({
            language: language,
            initComplete: initComplete,
            infoCallback: infoCallback,
        })

        ins.clear()
        ins.rows.add(response.data)
        ins.draw()

        // Добавление кнопок на таблицу
        const selBtn = document.getElementById('dt-length-0')
        const parentNode = selBtn.parentElement

        const deleteButton = document.createElement('input')
        deleteButton.type = 'button'
        deleteButton.value = 'Удалить'
        deleteButton.onclick = deleteRow

        const addButton = document.createElement('input')
        addButton.type = 'button'
        addButton.value = 'Добавить'
        addButton.onclick = addRow

        // Добавляем кнопки после кнопки выбора количества записей на странице
        parentNode.insertBefore(deleteButton, selBtn.nextSibling)
        parentNode.insertBefore(addButton, selBtn.nextSibling)
    })
}


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

// Получить таблицу
const getTable = async () => {
    try {
        const response = await axios.get('/api/get_table')
        return response
    } catch (error) {
        console.log(error)
    }
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

        // Если успешно загружено, обновляем данные в DataTable
        getTable().then((response) => {
            // Получаем ссылку на экземпляр DataTable
            const dataTable = $('#mainTable').DataTable()

            // Очищаем текущие данные в таблице
            dataTable.clear()

            // Добавляем новые данные в таблицу
            dataTable.rows.add(response.data)

            // Перерисовываем таблицу
            dataTable.draw()
        })
    }).catch((error) => {
        // Обрабатываем ошибку, если необходимо
        const text = 'Ошибка при загрузке файла'
        alert(text)
        console.error(text, error)
    })
}


$(document).ready(() => {
    init()
})
