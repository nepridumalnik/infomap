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

// Обработка загрузки страницы
const init = () => {
    getTable().then((response) => {
        const table = $('#mainTable').DataTable()

        table.clear()
        table.rows.add(response.data)
        table.draw()

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

// Завершение инициализации
const initComplete = () => {
    // Получаем данные для таблицы после инициализации
    getTable().then((response) => {
        const table = $('#mainTable').DataTable()
        table.clear()
        table.rows.add(response.data)
        table.draw()
    }).catch((error) => {
        console.error('Ошибка при загрузке данных для таблицы:', error)
    })
}

// Сделать таблицу редактируемой
const makeTableEditable = () => {
    // Находим таблицу по идентификатору
    const table = document.getElementById('mainTable')

    // Получаем все ячейки таблицы
    const cells = table.getElementsByTagName('td')

    // Добавляем обработчик двойного клика к каждой ячейке
    for (let i = 0; i < cells.length; i++) {
        cells[i].addEventListener('dblclick', function () {
            // Сохраняем текущее значение ячейки
            const currentValue = this.innerHTML

            // Создаем поле ввода
            let input = document.createElement('input')
            input.type = 'text'
            input.value = currentValue

            // Заменяем содержимое ячейки полем ввода
            this.innerHTML = ''
            this.appendChild(input)

            // Добавляем обработчик нажатия Enter для сохранения изменений
            input.addEventListener('keypress', function (event) {
                if (event.keyCode === 13) { // Enter key code
                    // При нажатии Enter сохраняем новое значение ячейки
                    let newValue = this.value
                    let parentCell = this.parentElement
                    parentCell.innerHTML = newValue
                }
            })

            // Добавляем обработчик события потери фокуса
            input.addEventListener('blur', function () {
                // Получаем все ячейки строки
                const rowCells = this.parentElement.parentElement.querySelectorAll('td')

                // Создаем массив для хранения содержимого ячеек
                const rowContentArray = []

                // Добавляем содержимое каждой ячейки в массив
                for (const cell of rowCells) {
                    let cellContent = cell.innerHTML
                    // Проверяем, является ли содержимое ячейки типом input
                    if (cell.querySelector('input')) {
                        // Если является, используем значение input
                        cellContent = cell.querySelector('input').value
                    }
                    rowContentArray.push(cellContent);
                }

                // Преобразуем массив в JSON-строку
                const rowContentJSON = JSON.stringify(rowContentArray)

                // Отправляем POST-запрос на сервер
                axios.put('/api/update_row', rowContentJSON)
                    .then(_ => {
                        console.log('Строка успешно обновлена')
                    })
                    .catch(error => {
                        console.error('Ошибка при обновлении строки:', error)
                    }).finally(() => {
                        // В любом случае обновляем таблицу
                        init()
                    })
            })

            // Фокусируемся на поле ввода
            input.focus()

            // Предотвращаем дальнейшее распространение события двойного клика
            event.stopPropagation()
        })
    }
}

// Обработка обновления информации в таблице
const infoCallback = (settings, start, end, max, total, pre) => {
    console.log('infoCallback()')
    makeTableEditable()

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
            }
        }
    })
}

// Функция для удаления строки
const deleteRow = () => {
    const selectedRow = $('#mainTable tbody .selected')
    if (selectedRow.length === 0) {
        const text = 'Строка не выделена'
        console.error(text)
        alert(text)
        return
    }

    // Получаем значение первого столбца из выделенной строки
    const id = selectedRow.find('td:first-child').text()

    // Отправляем DELETE запрос на /api/delete_row с параметром id
    axios.delete('/api/delete_row', {
        headers: {
            'Content-Type': 'application/json'
        },
        params: {
            id: id
        }
    }).then(_ => {
        console.log('Строка успешно удалена')
        const dataTable = $('#mainTable').DataTable()
        dataTable.row(selectedRow).remove().draw(false)
    }).catch(error => {
        console.error('Ошибка при удалении строки', error)
    })
}

// Функция для добавления строки
const addRow = () => {
    // Создание элементов формы
    let formContainer = document.createElement('div')
    formContainer.style.position = 'fixed'
    formContainer.style.top = '0'
    formContainer.style.left = '0'
    formContainer.style.width = '100%'
    formContainer.style.height = '100%'
    formContainer.style.backgroundColor = 'rgba(0, 0, 0, 0.5)'
    formContainer.style.display = 'flex'
    formContainer.style.alignItems = 'center'
    formContainer.style.justifyContent = 'center'
    formContainer.style.zIndex = '9999'

    let form = document.createElement('form')
    form.style.backgroundColor = '#fff'
    form.style.padding = '20px'
    form.style.borderRadius = '5px'
    form.style.display = 'flex'
    form.style.flexDirection = 'column' // Поля располагаются вертикально

    // Функция для создания поля ввода
    const createInputField = (labelText, inputType) => {
        let label = document.createElement('label')
        label.textContent = labelText + ': '
        let input = document.createElement('input')
        input.name = labelText
        input.type = inputType
        input.required = true
        input.style.marginBottom = '10px'
        label.appendChild(input)
        form.appendChild(label)
    }

    // Создание полей ввода
    createInputField('Регион', 'text')
    createInputField('Назначен ответственный', 'text')
    createInputField('Страница подтверждена', 'text')
    createInputField('ВКонтакте', 'url')
    createInputField('Одноклассники', 'url')
    createInputField('Telegram', 'url')
    createInputField('Официальная страница не ведется на основании', 'text')
    createInputField('Комментарий по НПА', 'text')
    createInputField('Полное наименование объекта', 'text')
    createInputField('ОГРН', 'text')
    createInputField('Статус', 'text')
    createInputField('Комментарий', 'textarea')

    // Создание кнопок
    let buttonsContainer = document.createElement('div')
    buttonsContainer.style.display = 'flex'
    buttonsContainer.style.justifyContent = 'space-between'
    form.appendChild(buttonsContainer)

    let submitButton = document.createElement('button')
    submitButton.type = 'submit'
    submitButton.textContent = 'Отправить'
    submitButton.style.padding = '10px 20px'
    submitButton.style.backgroundColor = '#4CAF50'
    submitButton.style.color = '#fff'
    submitButton.style.border = 'none'
    submitButton.style.borderRadius = '3px'
    submitButton.style.cursor = 'pointer'
    buttonsContainer.appendChild(submitButton)

    let cancelButton = document.createElement('button')
    cancelButton.type = 'button'
    cancelButton.textContent = 'Отмена'
    cancelButton.style.padding = '10px 20px'
    cancelButton.style.backgroundColor = '#f44336'
    cancelButton.style.color = '#fff'
    cancelButton.style.border = 'none'
    cancelButton.style.borderRadius = '3px'
    cancelButton.style.cursor = 'pointer'
    buttonsContainer.appendChild(cancelButton)

    // Добавление формы на страницу
    formContainer.appendChild(form)
    document.body.appendChild(formContainer)

    // Закрытие формы при нажатии на серый фон или кнопку 'Отмена'
    const closeForm = () => {
        document.body.removeChild(formContainer)
    }

    formContainer.addEventListener('click', (event) => {
        if (event.target === formContainer) {
            closeForm()
        }
    })

    cancelButton.addEventListener('click', closeForm)

    // Отправка формы
    form.addEventListener('submit', (event) => {
        const formData = new FormData(form)

        // Отправляем данные на сервер
        axios.post('/api/add_row', formData).then((response) => {
            console.log('response: ' + JSON.stringify(response.data))
            const row = []

            for (const value of Object.values(response.data)) {
                row.push(value)
            }

            const table = $('#mainTable').DataTable()
            table.row.add(row)
            table.draw()
            closeForm()
        }).catch((error) => {
            console.error(error)
            alert(error)
        })

        // Отменяем стандартное поведение формы
        event.preventDefault()
    })
}

// Деавторизация
const unauth = () => {
    console.log('unauth')
    axios.post('/unauth', {}).then((response) => {
        console.log('Response: ' + response)
        window.location.reload()
    }).catch((error) => {
        console.error('Error: ' + error)
    })
}

// Получить таблицу
const getTable = async () => {
    try {
        const response = await axios.get('/api/get_table')
        return response
    } catch (error) {
        console.error(error)
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
    // Первичная инициализация таблицы
    $('#mainTable').DataTable({
        language: language,
        initComplete: initComplete,
        infoCallback: infoCallback,
    })

    // Обработчики событий клика для кнопок
    $('#mainTable').on('click', 'input[value=\'Удалить\']', deleteRow)
    $('#mainTable').on('click', 'input[value=\'Добавить\']', addRow)

    init()
})
