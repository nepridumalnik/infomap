const Block = () => {
    return <h2>Это блок</h2>
}

class Button extends React.Component {
    render() {
        const increment = () => {
            this.props.onClick()
        }


        return (<button onClick={increment}> Кнопа</button >)
    }
}

class App extends React.Component {
    render() {
        let counter = 0
        const click = () => {
            counter++
            console.log("counter: " + counter)
        }

        return (
            <div>
                <h1>Я компонент App: {counter}</h1>
                <Block />
                <Button onClick={click} />
            </div>
        )
    }
}

ReactDOM.render(
    <App>

    </App>, document.getElementById("app")
)
