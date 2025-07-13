// ユーザーマスタ管理画面のJavaScript

let currentEditUserId = null;

// ページ読み込み時
document.addEventListener('DOMContentLoaded', function() {
    loadUsers();
    loadOrganizations();
});

// ユーザー一覧を読み込み
async function loadUsers() {
    try {
        const response = await fetch('/api/users');
        if (response.ok) {
            const users = await response.json();
            displayUsers(users);
        } else {
            console.error('ユーザー一覧の取得に失敗しました');
        }
    } catch (error) {
        console.error('エラー:', error);
    }
}

// ユーザー一覧を表示
function displayUsers(users) {
    const tbody = document.getElementById('userTableBody');
    tbody.innerHTML = '';
    
    users.forEach(user => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${user.user_id}</td>
            <td>${user.name}</td>
            <td>${getRoleName(user.role_id)}</td>
            <td>${getOrganizationNames(user.organizations)}</td>
            <td>${user.status === 'active' ? '有効' : '無効'}</td>
            <td>
                <button class="btn btn-sm" onclick="editUser('${user.user_id}')">編集</button>
                <button class="btn btn-sm btn-danger" onclick="deleteUser('${user.user_id}')">削除</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// ロール名を取得
function getRoleName(roleId) {
    const roles = {
        'R01': '管理職',
        'R02': '営業担当',
        'R03': '在庫管理者',
        'R04': 'システム管理者'
    };
    return roles[roleId] || roleId;
}

// 組織名を取得
function getOrganizationNames(organizations) {
    if (!organizations || organizations.length === 0) {
        return '未所属';
    }
    return organizations.map(org => org.name).join(', ');
}

// 組織一覧を読み込み
async function loadOrganizations() {
    try {
        const response = await fetch('/api/organizations');
        if (response.ok) {
            const organizations = await response.json();
            window.organizationList = organizations;
        }
    } catch (error) {
        console.error('組織一覧の取得に失敗:', error);
    }
}

// 新規ユーザー追加モーダルを表示
function showAddUserModal() {
    currentEditUserId = null;
    document.getElementById('modalTitle').textContent = '新規ユーザー追加';
    document.getElementById('userId').value = '';
    document.getElementById('userId').removeAttribute('readonly');
    document.getElementById('userName').value = '';
    document.getElementById('userRole').value = '';
    document.getElementById('userStatus').value = 'active';
    
    // 組織選択を表示
    displayOrganizationCheckboxes([]);
    
    document.getElementById('userModal').style.display = 'block';
}

// ユーザー編集
async function editUser(userId) {
    currentEditUserId = userId;
    
    try {
        const response = await fetch(`/api/users/${userId}`);
        if (response.ok) {
            const user = await response.json();
            
            document.getElementById('modalTitle').textContent = 'ユーザー編集';
            document.getElementById('userId').value = user.user_id;
            document.getElementById('userId').setAttribute('readonly', true);
            document.getElementById('userName').value = user.name;
            document.getElementById('userRole').value = user.role_id;
            document.getElementById('userStatus').value = user.status || 'active';
            
            // 組織選択を表示
            displayOrganizationCheckboxes(user.organizations || []);
            
            document.getElementById('userModal').style.display = 'block';
        }
    } catch (error) {
        console.error('ユーザー情報の取得に失敗:', error);
        alert('ユーザー情報の取得に失敗しました');
    }
}

// 組織選択チェックボックスを表示
function displayOrganizationCheckboxes(userOrganizations) {
    const container = document.getElementById('userOrgContainer');
    container.innerHTML = '';
    
    if (!window.organizationList || window.organizationList.length === 0) {
        container.innerHTML = '<p>組織データがありません</p>';
        return;
    }
    
    const userOrgIds = userOrganizations.map(org => org.org_id);
    
    window.organizationList.forEach(org => {
        const div = document.createElement('div');
        div.style.marginBottom = '5px';
        
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.id = `org_${org.org_id}`;
        checkbox.value = org.org_id;
        checkbox.checked = userOrgIds.includes(org.org_id);
        
        const label = document.createElement('label');
        label.htmlFor = `org_${org.org_id}`;
        label.textContent = ` ${org.name} (${org.org_id})`;
        label.style.marginLeft = '5px';
        
        div.appendChild(checkbox);
        div.appendChild(label);
        container.appendChild(div);
    });
}

// モーダルを閉じる
function closeUserModal() {
    document.getElementById('userModal').style.display = 'none';
    currentEditUserId = null;
}

// フォーム送信
document.getElementById('userForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const userId = document.getElementById('userId').value;
    const userName = document.getElementById('userName').value;
    const roleId = document.getElementById('userRole').value;
    const status = document.getElementById('userStatus').value;
    
    // 選択された組織を取得
    const selectedOrgs = [];
    const checkboxes = document.querySelectorAll('#userOrgContainer input[type="checkbox"]:checked');
    checkboxes.forEach(checkbox => {
        selectedOrgs.push(checkbox.value);
    });
    
    const userData = {
        user_id: userId,
        name: userName,
        role_id: roleId,
        status: status,
        organizations: selectedOrgs
    };
    
    try {
        let response;
        if (currentEditUserId) {
            // 更新
            response = await fetch(`/api/users/${currentEditUserId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
        } else {
            // 新規作成
            response = await fetch('/api/users', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
        }
        
        if (response.ok) {
            alert(currentEditUserId ? 'ユーザー情報を更新しました' : 'ユーザーを追加しました');
            closeUserModal();
            loadUsers();
        } else {
            const error = await response.json();
            alert('エラー: ' + (error.message || '保存に失敗しました'));
        }
    } catch (error) {
        console.error('エラー:', error);
        alert('保存中にエラーが発生しました');
    }
});

// ユーザー削除
async function deleteUser(userId) {
    if (!confirm(`ユーザー ${userId} を削除してもよろしいですか？`)) {
        return;
    }
    
    try {
        const response = await fetch(`/api/users/${userId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            alert('ユーザーを削除しました');
            loadUsers();
        } else {
            alert('削除に失敗しました');
        }
    } catch (error) {
        console.error('エラー:', error);
        alert('削除中にエラーが発生しました');
    }
}

// モーダル外クリックで閉じる
window.onclick = function(event) {
    if (event.target == document.getElementById('userModal')) {
        closeUserModal();
    }
}