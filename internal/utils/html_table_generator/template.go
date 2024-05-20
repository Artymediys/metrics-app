package html_table_generator

const htmlTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
<style>
    table {
        width: 100%;
        border: 2px solid black;
        border-collapse: collapse;
    }
    th {
        text-align: left;
        background: #ccc;
        padding: 5px;
        border: 1px solid black;
    }
    td {
        padding: 5px;
        border: 1px solid black;
    }
</style>
</head>
<body>
    <table>
        <tr>
            {{range .Headers}}<th>{{.}}</th>
            {{end}}
        </tr>
        {{range .Rows}}<tr>
            {{range .}}<td>{{.}}</td>
            {{end}}
        </tr>
        {{end}}
    </table>
</body>
</html>
`
