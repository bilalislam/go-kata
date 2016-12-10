var userGridRow = React.createClass({
    render: function () {
        return (
            <tr>
                <td>{this.props.item.Name}</td>
                <td>{this.props.item.Age}</td>
                <td>{this.props.item.Status}</td>
            </tr>
        )
    }
});

var userGridTable = React.createClass({
    getInitialState: function () {
        return {
            items: []
        }
    },
    componentDidMount: function () {
        $.get(this.props.url, function (data) {
            if (this.isMounted) {
                this.setState({
                    items: data
                })
            }
        }.bind(this))
    },
    render: function () {
        var rows = [];
        var GridRowFactory = React.createFactory(userGridRow);
        this.state.items.forEach(function (item) {
            rows.push(GridRowFactory({item: item}));
        });

        return (
            <table className="table table-bordered table-responsive">
                <thead>
                <tr>
                    <th>Name</th>
                    <th>Age</th>
                    <th>State</th>
                </tr>
                </thead>
                <tbody>
                {rows}
                </tbody>
            </table>
        )
    }
});

var ExampleApplicationFactory = React.createFactory(userGridTable);

ReactDOM.render(
    ExampleApplicationFactory({url: "/GetAllUsers"}),
    document.getElementById('griddata')
);

