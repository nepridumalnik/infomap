// https://react.dev/learn
// https://ru.react.js.org/docs/react-component.html

class Button extends React.Component {
    render() {
        const increment = () => {
            this.props.onClick()
        }

        return (<button onClick={increment}> Кнопа</button >)
    }
}

class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = { counter: 0 }
    }

    render() {
        const increment = () => {
            this.setState({ counter: this.state.counter + 1 })
        }

        return (
            <div>
                <div>Было сделано "{this.state.counter}" нажатий</div>
                <Button onClick={increment} />
            </div>
        )
    }
}

ReactDOM.render(
    <App>

    </App>, document.getElementById("app")
)
