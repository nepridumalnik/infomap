class Button extends React.Component {
    render() {
        const increment = () => {
            this.props.onClick()
        }

        return (<button onClick={increment}>Кнопка</button >)
    }
}

export { Button }
