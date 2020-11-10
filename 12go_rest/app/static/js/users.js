const usersModule = (() => {
    const BASE_URL = "http://localhost:58606/api/v1/users";

    const headers = new Headers();
    headers.set("Content-type", "application/json");

    const handleError = async(res) => {
        const resJson = await res.json();

        switch (res.status) {
            case 200:
                alert(resJson.message);
                window.location.href = "/";
                break;
            case 201:
                alert(resJson.message);
                window.location.href = "/";
                break;
            case 400:
                alert(resJson.error);
                break;
            case 404:
                alert(resJson.error);
                break;
            case 500:
                alert(resJson.error);
                break;
            default:
                alert("予期せぬエラー");
        }
    }

    return {
        fetchAllUsers: async() => {
            const res = await fetch(BASE_URL);
            const users = await res.json();

            for (let i=0; i < users.length; i++) {
                const user = users[i];
                const body = `
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
                document.getElementById("users-list").insertAdjacentHTML("beforeend", body);
            }
        },
        createUser: async() => {
            const name = document.getElementById("name").value;
            const profile = document.getElementById("profile").value;
            const dateOfBirth = document.getElementById("date-of-birth").value;

            const body = {
                name: name,
                profile: profile,
                date_of_birth: dateOfBirth
            }

            const res = await fetch(BASE_URL, {
                method: "POST",
                headers: headers,
                body: JSON.stringify(body)
            });

            return handleError(res);
        },
        setExistingValue: async(uid) => {
            const res = await fetch(BASE_URL + "/" + uid);
            const resJson = await res.json();
            console.log(resJson);

            document.getElementById("name").value = resJson.name;
            document.getElementById("profile").value = resJson.profile;
            document.getElementById("date-of-birth").value = resJson.date_of_birth;
        },
        saveUser: async(uid) => {
            const name = document.getElementById("name").value;
            const profile = document.getElementById("profile").value;
            const dateOfBirth = document.getElementById("date-of-birth").value;

            const body = {
                name: name,
                profile: profile,
                date_of_birth: dateOfBirth
            }

            const res = await fetch(BASE_URL + "/" + uid, {
                method: "PUT",
                headers: headers,
                body: JSON.stringify(body)
            });
            
            return handleError(res)
        },
        deleteUser: async(uid) => {
            const ret = window.confirm("このユーザを削除しますか？");
            
            if (!ret) {
                return false
            } else {
                const res = await fetch(BASE_URL + "/" + uid, {
                    method: "DELETE",
                    headers: headers
                });

                return handleError(res);
            }
        }
    }
})();