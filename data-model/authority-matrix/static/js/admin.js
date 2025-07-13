// API呼び出しのベースURL（セッションベース）
function apiUrl(path) {
    // セッションベースなのでuser_idパラメータは不要
    return path;
}

// タブ切り替え
function showTab(tabId, element) {
    // すべてのタブコンテンツを非表示
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // すべてのタブリンクから active クラスを削除
    document.querySelectorAll('.nav-tabs a').forEach(link => {
        link.classList.remove('active');
    });
    
    // 選択されたタブを表示
    document.getElementById(tabId).classList.add('active');
    element.classList.add('active');
}

// アラート表示
function showAlert(message, type = 'info') {
    const alertContainer = document.getElementById('alert-container');
    const alert = document.createElement('div');
    alert.className = `alert alert-${type}`;
    alert.textContent = message;
    
    alertContainer.appendChild(alert);
    
    // 3秒後に自動的に削除
    setTimeout(() => {
        alert.remove();
    }, 3000);
}

// モーダル表示/非表示
function showModal(modalId) {
    document.getElementById(modalId).style.display = 'block';
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

// メニュー設定関連
async function loadMenuSettings() {
    try {
        const response = await fetch(apiUrl('/api/menu-settings'));
        const data = await response.json();
        
        const tbody = document.querySelector('#menu-table tbody');
        tbody.innerHTML = '';
        
        data.forEach(menu => {
            const row = tbody.insertRow();
            row.innerHTML = `
                <td>${menu.menu_id}</td>
                <td>${menu.menu_name}</td>
                <td>${menu.screen_id}</td>
                <td>
                    <button class="btn btn-sm" onclick="editMenu('${menu.menu_id}', '${menu.menu_name}', '${menu.screen_id}')">編集</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteMenu('${menu.menu_id}')">削除</button>
                </td>
            `;
        });
    } catch (error) {
        showAlert('メニュー設定の読み込みに失敗しました', 'error');
    }
}

function showAddMenuModal() {
    document.getElementById('menu-modal-title').textContent = '新規メニュー追加';
    document.getElementById('menu-form').reset();
    document.getElementById('menu-id-original').value = '';
    showModal('menu-modal');
}

function editMenu(menuId, menuName, screenId) {
    document.getElementById('menu-modal-title').textContent = 'メニュー編集';
    document.getElementById('menu-id-original').value = menuId;
    document.getElementById('menu-id').value = menuId;
    document.getElementById('menu-name').value = menuName;
    document.getElementById('screen-id').value = screenId;
    showModal('menu-modal');
}

async function saveMenu(event) {
    event.preventDefault();
    
    const originalId = document.getElementById('menu-id-original').value;
    const menuData = {
        menu_id: document.getElementById('menu-id').value,
        menu_name: document.getElementById('menu-name').value,
        screen_id: document.getElementById('screen-id').value
    };
    
    try {
        const method = originalId ? 'PUT' : 'POST';
        const url = originalId ? `/api/menu-settings/${originalId}` : '/api/menu-settings';
        
        const response = await fetch(apiUrl(url), {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(menuData)
        });
        
        if (response.ok) {
            showAlert('メニューを保存しました', 'success');
            closeModal('menu-modal');
            loadMenuSettings();
        } else {
            showAlert('保存に失敗しました', 'error');
        }
    } catch (error) {
        showAlert('エラーが発生しました', 'error');
    }
}

async function deleteMenu(menuId) {
    if (!confirm(`メニューID: ${menuId} を削除しますか？`)) {
        return;
    }
    
    try {
        const response = await fetch(apiUrl(`/api/menu-settings/${menuId}`), {
            method: 'DELETE'
        });
        
        if (response.ok) {
            showAlert('メニューを削除しました', 'success');
            loadMenuSettings();
        } else {
            showAlert('削除に失敗しました', 'error');
        }
    } catch (error) {
        showAlert('エラーが発生しました', 'error');
    }
}

// APIマッピング関連
async function loadApiMappings() {
    try {
        const response = await fetch(apiUrl('/api/api-screen-mappings'));
        const data = await response.json();
        
        const tbody = document.querySelector('#mapping-table tbody');
        tbody.innerHTML = '';
        
        data.forEach(mapping => {
            const row = tbody.insertRow();
            row.innerHTML = `
                <td>${mapping.screen_id}</td>
                <td>${mapping.api_id}</td>
                <td>${mapping.description || ''}</td>
                <td>
                    <button class="btn btn-sm" onclick="editMapping('${mapping.screen_id}', '${mapping.api_id}', '${mapping.description || ''}')">編集</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteMapping('${mapping.screen_id}', '${mapping.api_id}')">削除</button>
                </td>
            `;
        });
    } catch (error) {
        showAlert('APIマッピングの読み込みに失敗しました', 'error');
    }
}

function showAddMappingModal() {
    document.getElementById('mapping-modal-title').textContent = '新規マッピング追加';
    document.getElementById('mapping-form').reset();
    document.getElementById('mapping-original-screen-id').value = '';
    document.getElementById('mapping-original-api-id').value = '';
    showModal('mapping-modal');
}

function editMapping(screenId, apiId, description) {
    document.getElementById('mapping-modal-title').textContent = 'マッピング編集';
    document.getElementById('mapping-original-screen-id').value = screenId;
    document.getElementById('mapping-original-api-id').value = apiId;
    document.getElementById('mapping-screen-id').value = screenId;
    document.getElementById('mapping-api-id').value = apiId;
    document.getElementById('mapping-description').value = description;
    showModal('mapping-modal');
}

async function saveMapping(event) {
    event.preventDefault();
    
    const originalScreenId = document.getElementById('mapping-original-screen-id').value;
    const originalApiId = document.getElementById('mapping-original-api-id').value;
    const mappingData = {
        screen_id: document.getElementById('mapping-screen-id').value,
        api_id: document.getElementById('mapping-api-id').value,
        description: document.getElementById('mapping-description').value
    };
    
    try {
        const method = originalScreenId ? 'PUT' : 'POST';
        const url = originalScreenId ? 
            `/api/api-screen-mappings/${originalScreenId}/${originalApiId}` : 
            '/api/api-screen-mappings';
        
        const response = await fetch(apiUrl(url), {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(mappingData)
        });
        
        if (response.ok) {
            showAlert('マッピングを保存しました', 'success');
            closeModal('mapping-modal');
            loadApiMappings();
        } else {
            showAlert('保存に失敗しました', 'error');
        }
    } catch (error) {
        showAlert('エラーが発生しました', 'error');
    }
}

async function deleteMapping(screenId, apiId) {
    if (!confirm(`画面ID: ${screenId}, API ID: ${apiId} のマッピングを削除しますか？`)) {
        return;
    }
    
    try {
        const response = await fetch(apiUrl(`/api/api-screen-mappings/${screenId}/${apiId}`), {
            method: 'DELETE'
        });
        
        if (response.ok) {
            showAlert('マッピングを削除しました', 'success');
            loadApiMappings();
        } else {
            showAlert('削除に失敗しました', 'error');
        }
    } catch (error) {
        showAlert('エラーが発生しました', 'error');
    }
}

// ロール権限関連
async function loadRoles() {
    try {
        const response = await fetch(apiUrl('/api/roles'));
        const data = await response.json();
        
        const select = document.getElementById('role-select');
        select.innerHTML = '<option value="">ロールを選択してください</option>';
        
        data.forEach(role => {
            const option = document.createElement('option');
            option.value = role.role_id;
            option.textContent = `${role.role_id} - ${role.role_name}`;
            select.appendChild(option);
        });
    } catch (error) {
        showAlert('ロールの読み込みに失敗しました', 'error');
    }
}

async function loadRolePermissions() {
    const roleId = document.getElementById('role-select').value;
    if (!roleId) {
        document.getElementById('role-permissions').innerHTML = '';
        return;
    }
    
    try {
        // 画面権限を取得
        const [screensResponse, roleScreensResponse] = await Promise.all([
            fetch(apiUrl('/api/screens')),
            fetch(apiUrl(`/api/role-screen-permissions/${roleId}`))
        ]);
        
        const screens = await screensResponse.json();
        const roleScreens = await roleScreensResponse.json();
        
        // 権限がある画面のIDとレベルをMapに格納
        const roleScreenPermissions = new Map();
        roleScreens.forEach(rs => {
            roleScreenPermissions.set(rs.screen_id, rs.permission_level || 'VIEWER');
        });
        
        let html = '<h4>画面権限</h4>';
        html += '<p class="info">画面へのアクセス権限を設定します。画面に紐づくAPIは自動的にアクセス可能になります。</p>';
        html += '<table><thead><tr><th>アクセス権限</th><th>画面ID</th><th>画面名</th><th>権限レベル</th></tr></thead><tbody>';
        
        screens.forEach(screen => {
            const hasPermission = roleScreenPermissions.has(screen.screen_id);
            const permissionLevel = roleScreenPermissions.get(screen.screen_id) || 'VIEWER';
            const checked = hasPermission ? 'checked' : '';
            const selectDisabled = hasPermission ? '' : 'disabled';
            
            html += `
                <tr>
                    <td><input type="checkbox" ${checked} onchange="updateScreenPermission('${roleId}', '${screen.screen_id}', this.checked, document.getElementById('screen-level-${screen.screen_id}').value)"></td>
                    <td>${screen.screen_id}</td>
                    <td>${screen.screen_name}</td>
                    <td>
                        <select id="screen-level-${screen.screen_id}" ${selectDisabled} onchange="updateScreenPermissionLevel('${roleId}', '${screen.screen_id}', this.value)">
                            <option value="VIEWER" ${permissionLevel === 'VIEWER' ? 'selected' : ''}>Viewer</option>
                            <option value="EDITOR" ${permissionLevel === 'EDITOR' ? 'selected' : ''}>Editor</option>
                        </select>
                    </td>
                </tr>
            `;
        });
        
        html += '</tbody></table>';
        document.getElementById('role-permissions').innerHTML = html;
        
    } catch (error) {
        showAlert('権限情報の読み込みに失敗しました', 'error');
    }
}

async function updateScreenPermission(roleId, screenId, hasPermission, permissionLevel = 'VIEWER') {
    try {
        const method = hasPermission ? 'POST' : 'DELETE';
        const body = hasPermission ? JSON.stringify({ permission_level: permissionLevel }) : undefined;
        
        const response = await fetch(apiUrl(`/api/role-screen-permissions/${roleId}/${screenId}`), {
            method: method,
            headers: hasPermission ? { 'Content-Type': 'application/json' } : {},
            body: body
        });
        
        if (response.ok) {
            showAlert(`画面権限を${hasPermission ? '付与' : '削除'}しました`, 'success');
            // 権限レベルセレクトボックスを有効/無効化
            const select = document.getElementById(`screen-level-${screenId}`);
            if (select) {
                select.disabled = !hasPermission;
            }
        } else {
            showAlert('権限の更新に失敗しました', 'error');
        }
    } catch (error) {
        showAlert('エラーが発生しました', 'error');
    }
}

async function updateScreenPermissionLevel(roleId, screenId, permissionLevel) {
    try {
        const response = await fetch(apiUrl(`/api/role-screen-permissions/${roleId}/${screenId}`), {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ permission_level: permissionLevel })
        });
        
        if (response.ok) {
            showAlert('画面権限レベルを更新しました', 'success');
        } else {
            showAlert('権限レベルの更新に失敗しました', 'error');
        }
    } catch (error) {
        showAlert('エラーが発生しました', 'error');
    }
}


// ページ読み込み時の初期化
document.addEventListener('DOMContentLoaded', () => {
    loadMenuSettings();
    loadApiMappings();
    loadRoles();
});

// モーダルの外側クリックで閉じる
window.onclick = function(event) {
    if (event.target.classList.contains('modal')) {
        event.target.style.display = 'none';
    }
}