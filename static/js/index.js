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

// Завершение обновления таблицы
const infoCallback = (settings, start, end, max, total, pre) => {
    const tbody = document.querySelector('#mainTable tbody')
    const rows = tbody.querySelectorAll('tr')

    // Обработчик события клика для каждой строки таблицы
    rows.forEach((row) => {
        row.addEventListener('click', () => {
            // Сбрасываем выделение предыдущей выделенной строки
            const previouslySelected = tbody.querySelector('.selected')
            if (previouslySelected) {
                previouslySelected.classList.remove('selected')
            }
            // Выделяем текущую строку
            row.classList.add('selected')
            // Выводим содержимое выделяемой строки в консоль
            console.log("Содержимое выделенной строки:")
            row.querySelectorAll('td').forEach((cell) => {
                console.log(cell.textContent)
            })
        })
    })

    // Обработчик события потери фокуса для всего документа
    document.addEventListener('click', (event) => {
        const isClickedOutsideTable = !event.target.closest('#mainTable')
        if (isClickedOutsideTable) {
            // Сбрасываем выделение при клике вне таблицы
            const selectedRow = tbody.querySelector('.selected')
            if (selectedRow) {
                selectedRow.classList.remove('selected')
                // Выводим сообщение о потере фокуса в консоль
                console.log("Фокус потерян")
            }
        }
    })

    // Остальной ваш код
    rows.forEach((row) => {
        const cells = row.querySelectorAll('td')
        cells.forEach((cell) => {
            console.log(cell.textContent)
        })
    })
}

// Функция для удаления строки
const deleteRow = () => {
    const selectedRow = document.querySelector('#mainTable tbody .selected')
    if (!selectedRow) {
        console.log('Строка не выделена')
        alert('Строка не выделена')
        return
    }

    // Получаем значение первого столбца из выделенной строки
    const id = selectedRow.querySelector('td:first-child').textContent

    // Отправляем DELETE запрос на /api/delete_row с параметром firstColumnValue
    axios.delete('/api/delete_row', {
        headers: {
            'Content-Type': 'application/json'
        },
        params: {
            id: id
        }
    }).then(response => {
        console.log('Строка успешно удалена')
    }).catch(error => {
        console.error('Ошибка при удалении строки', error)
    })
}

// Функция для добавления строки
const addRow = () => {
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
    // Обработчики событий клика для кнопок
    $('#mainTable').on('click', 'input[value="Удалить"]', deleteRow)
    $('#mainTable').on('click', 'input[value="Добавить"]', addRow)

    init()
})
