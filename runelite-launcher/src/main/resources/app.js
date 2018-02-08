var h = picodom.h

function UI(c) {
    return h('div', { class: "container", style: { padding: "25px" } },
        h('div', { class: "progress", style: { height: "35px", marginBottom: "15px" }},
            h('div', { class: "progress-bar progress-bar-striped progress-bar-animated", style: { width: c.data.progress + "%"}}, c.data.progress)),
        h('ul', { class: "list-group"},
            c.data.label.split("\n").filter(function (m) { return m !== "" }).reverse().map(function (message) {
                return h('li', { class: "list-group-item" }, message)
            })))
}

var node;
function render() {
    node = picodom.patch(node, node=UI(controller))
}

controller.render = render;
render();