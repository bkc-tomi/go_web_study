const searchModule = (() => {
    const BASE_URL = "http://localhost:58606/api/v1/search";

    return {
        searchUsers: async() => {
            const query = document.getElementById("search");
            const res = await fetch(BASE_URL + "?q=" + query.value);
            const result = await res.json();

            let body = "";
            
            for (let i=0; i<result.length; i++) {
                const user = result[i];
                body += `
                <tr>
                    <td>${user.id}</td>
                    <td>${user.name}</td>
                    <td>${user.profile}</td>
                    <td>${user.date_of_birth}</td>
                    <td>${user.create_at}</td>
                    <td>${user.update_at}</td>
                    <td><a href="edit?uid=${user.id}">編集</a></td>
                </tr>
                `;
            }
            document.getElementById("users-list").innerHTML = body;
        }
    }
})();