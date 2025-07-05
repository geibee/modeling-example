const API_BASE = '/api';

let departments = [];
let attributes = [];

// 初期化
document.addEventListener('DOMContentLoaded', () => {
    loadDepartments();
    loadAttributes();
});

// タブ切り替え
function showTab(tabName) {
    const tabs = document.querySelectorAll('.tab-content');
    const buttons = document.querySelectorAll('.tab-button');
    
    tabs.forEach(tab => {
        tab.classList.remove('active');
        if (tab.id === tabName) {
            tab.classList.add('active');
        }
    });
    
    buttons.forEach(button => {
        button.classList.remove('active');
        if (button.textContent.includes(tabName === 'departments' ? '部門管理' : '期間別組織属性')) {
            button.classList.add('active');
        }
    });
}

// 部門管理
async function loadDepartments() {
    try {
        const response = await fetch(`${API_BASE}/departments`);
        departments = await response.json();
        renderDepartments();
        updateDepartmentSelects();
    } catch (error) {
        console.error('部門の読み込みに失敗しました:', error);
    }
}

function renderDepartments() {
    const tbody = document.getElementById('departmentsList');
    tbody.innerHTML = departments.map(dept => `
        <tr>
            <td>${dept.department_id}</td>
            <td>${dept.department_name}</td>
            <td>${formatDateTime(dept.created_at)}</td>
            <td>${formatDateTime(dept.updated_at)}</td>
            <td>
                <button class="edit" onclick="editDepartment('${dept.department_id}')">編集</button>
                <button class="delete" onclick="deleteDepartment('${dept.department_id}')">削除</button>
            </td>
        </tr>
    `).join('');
}

function showDepartmentForm() {
    document.getElementById('deptEditMode').value = 'create';
    document.getElementById('deptId').value = '';
    document.getElementById('deptId').disabled = false;
    document.getElementById('deptName').value = '';
    document.getElementById('departmentForm').style.display = 'block';
}

function hideDepartmentForm() {
    document.getElementById('departmentForm').style.display = 'none';
}

async function saveDepartment(event) {
    event.preventDefault();
    
    const mode = document.getElementById('deptEditMode').value;
    const deptId = document.getElementById('deptId').value;
    const deptName = document.getElementById('deptName').value;
    
    const data = {
        department_id: deptId,
        department_name: deptName
    };
    
    try {
        let response;
        if (mode === 'create') {
            response = await fetch(`${API_BASE}/departments`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
        } else {
            response = await fetch(`${API_BASE}/departments/${deptId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ department_name: deptName })
            });
        }
        
        if (response.ok) {
            hideDepartmentForm();
            loadDepartments();
            loadAttributes();
        } else {
            const error = await response.json();
            alert('エラー: ' + error.error);
        }
    } catch (error) {
        alert('保存に失敗しました: ' + error.message);
    }
}

async function editDepartment(deptId) {
    const dept = departments.find(d => d.department_id === deptId);
    if (!dept) return;
    
    document.getElementById('deptEditMode').value = 'edit';
    document.getElementById('deptId').value = dept.department_id;
    document.getElementById('deptId').disabled = true;
    document.getElementById('deptName').value = dept.department_name;
    document.getElementById('departmentForm').style.display = 'block';
}

async function deleteDepartment(deptId) {
    if (!confirm(`部門 ${deptId} を削除しますか？`)) return;
    
    try {
        const response = await fetch(`${API_BASE}/departments/${deptId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            loadDepartments();
            loadAttributes();
        } else {
            const error = await response.json();
            alert('エラー: ' + error.error);
        }
    } catch (error) {
        alert('削除に失敗しました: ' + error.message);
    }
}

// 期間別組織属性
async function loadAttributes() {
    try {
        const response = await fetch(`${API_BASE}/organization-attributes`);
        attributes = await response.json();
        renderAttributes();
    } catch (error) {
        console.error('組織属性の読み込みに失敗しました:', error);
    }
}

function renderAttributes(filteredAttrs = null) {
    const tbody = document.getElementById('attributesList');
    const attrsToRender = filteredAttrs || attributes;
    
    tbody.innerHTML = attrsToRender.map(attr => {
        const dept = departments.find(d => d.department_id === attr.department_id);
        const parentDept = attr.parent_department_id ? 
            departments.find(d => d.department_id === attr.parent_department_id) : null;
        
        return `
            <tr>
                <td>${attr.department_id}</td>
                <td>${dept ? dept.department_name : ''}</td>
                <td>${formatDate(attr.effective_date)}</td>
                <td>${attr.expiration_date ? formatDate(attr.expiration_date) : '現在'}</td>
                <td>${attr.parent_department_id || ''}</td>
                <td>${parentDept ? parentDept.department_name : ''}</td>
                <td>
                    <button class="edit" onclick="editAttribute('${attr.department_id}', '${formatDate(attr.effective_date)}')">編集</button>
                    <button class="delete" onclick="deleteAttribute('${attr.department_id}', '${formatDate(attr.effective_date)}')">削除</button>
                </td>
            </tr>
        `;
    }).join('');
}

function updateDepartmentSelects() {
    const deptFilter = document.getElementById('deptFilter');
    const attrDeptId = document.getElementById('attrDeptId');
    const attrParentDeptId = document.getElementById('attrParentDeptId');
    
    const deptOptions = departments.map(dept => 
        `<option value="${dept.department_id}">${dept.department_id} - ${dept.department_name}</option>`
    ).join('');
    
    deptFilter.innerHTML = '<option value="">全部門</option>' + deptOptions;
    attrDeptId.innerHTML = deptOptions;
    attrParentDeptId.innerHTML = '<option value="">なし（最上位）</option>' + deptOptions;
}

function filterAttributes() {
    const deptId = document.getElementById('deptFilter').value;
    if (deptId) {
        const filtered = attributes.filter(attr => attr.department_id === deptId);
        renderAttributes(filtered);
    } else {
        renderAttributes();
    }
}

function showAttributeForm() {
    document.getElementById('attrEditMode').value = 'create';
    document.getElementById('attrDeptId').disabled = false;
    document.getElementById('attrDeptId').value = '';
    document.getElementById('attrEffectiveDate').disabled = false;
    document.getElementById('attrEffectiveDate').value = '';
    document.getElementById('attrParentDeptId').value = '';
    document.getElementById('attributeForm').style.display = 'block';
}

function hideAttributeForm() {
    document.getElementById('attributeForm').style.display = 'none';
}

async function saveAttribute(event) {
    event.preventDefault();
    
    const mode = document.getElementById('attrEditMode').value;
    const deptId = document.getElementById('attrDeptId').value;
    const effectiveDate = document.getElementById('attrEffectiveDate').value;
    const parentDeptId = document.getElementById('attrParentDeptId').value || null;
    
    const data = {
        department_id: deptId,
        effective_date: effectiveDate,
        parent_department_id: parentDeptId
    };
    
    try {
        let response;
        if (mode === 'create') {
            response = await fetch(`${API_BASE}/organization-attributes`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
        } else {
            const originalDeptId = document.getElementById('attrOriginalDeptId').value;
            const originalDate = document.getElementById('attrOriginalEffectiveDate').value;
            response = await fetch(`${API_BASE}/organization-attributes/${originalDeptId}/${originalDate}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ parent_department_id: parentDeptId })
            });
        }
        
        if (response.ok) {
            hideAttributeForm();
            loadAttributes();
        } else {
            const error = await response.json();
            alert('エラー: ' + error.error);
        }
    } catch (error) {
        alert('保存に失敗しました: ' + error.message);
    }
}

async function editAttribute(deptId, effectiveDate) {
    const attr = attributes.find(a => 
        a.department_id === deptId && formatDate(a.effective_date) === effectiveDate
    );
    if (!attr) return;
    
    document.getElementById('attrEditMode').value = 'edit';
    document.getElementById('attrOriginalDeptId').value = deptId;
    document.getElementById('attrOriginalEffectiveDate').value = effectiveDate;
    document.getElementById('attrDeptId').value = attr.department_id;
    document.getElementById('attrDeptId').disabled = true;
    document.getElementById('attrEffectiveDate').value = effectiveDate;
    document.getElementById('attrEffectiveDate').disabled = true;
    document.getElementById('attrParentDeptId').value = attr.parent_department_id || '';
    document.getElementById('attributeForm').style.display = 'block';
}

async function deleteAttribute(deptId, effectiveDate) {
    if (!confirm(`部門 ${deptId} の ${effectiveDate} の属性を削除しますか？`)) return;
    
    try {
        const response = await fetch(`${API_BASE}/organization-attributes/${deptId}/${effectiveDate}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            loadAttributes();
        } else {
            const error = await response.json();
            alert('エラー: ' + error.error);
        }
    } catch (error) {
        alert('削除に失敗しました: ' + error.message);
    }
}

// ユーティリティ関数
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
}

function formatDateTime(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ja-JP');
}