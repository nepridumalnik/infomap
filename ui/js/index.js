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
