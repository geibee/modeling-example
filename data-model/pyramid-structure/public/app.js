const API_BASE = '/api';

let departments = [];
let attributes = [];

// 初期化
document.addEventListener('DOMContentLoaded', () => {
    loadDepartments();
    loadAttributes();
    
    // 階層ビューの日付を今日に設定
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('hierarchyDate').value = today;
    document.getElementById('reorganizeDate').value = today;
    document.getElementById('moveEffectiveDate').value = today;
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
        const buttonText = button.textContent;
        if ((tabName === 'departments' && buttonText.includes('部門管理')) ||
            (tabName === 'attributes' && buttonText.includes('期間別組織属性')) ||
            (tabName === 'hierarchy' && buttonText.includes('組織階層')) ||
            (tabName === 'reorganize' && buttonText.includes('組織再編'))) {
            button.classList.add('active');
        }
    });
    
    // 組織階層タブが選択されたら階層データを読み込む
    if (tabName === 'hierarchy') {
        updateHierarchy();
    } else if (tabName === 'reorganize') {
        updateReorganizeView();
        updateMoveSelects();
    }
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

// 組織階層ビュー
let hierarchyData = null;

async function updateHierarchy() {
    const date = document.getElementById('hierarchyDate').value;
    if (!date) return;
    
    const container = document.getElementById('hierarchyView');
    container.innerHTML = '<div class="hierarchy-loading">データを読み込み中...</div>';
    
    try {
        const response = await fetch(`${API_BASE}/hierarchy?date=${date}`);
        const data = await response.json();
        
        hierarchyData = buildHierarchyTree(data.attributes);
        renderHierarchy(hierarchyData, container);
    } catch (error) {
        console.error('階層データの読み込みに失敗しました:', error);
        container.innerHTML = '<div class="error">階層データの読み込みに失敗しました</div>';
    }
}

function buildHierarchyTree(attrs) {
    const nodeMap = new Map();
    const roots = [];
    
    // すべてのノードを作成
    attrs.forEach(attr => {
        const dept = departments.find(d => d.department_id === attr.department_id);
        nodeMap.set(attr.department_id, {
            id: attr.department_id,
            name: dept ? dept.department_name : attr.department_id,
            effectiveDate: formatDate(attr.effective_date),
            expirationDate: attr.expiration_date ? formatDate(attr.expiration_date) : null,
            parentId: attr.parent_department_id,
            children: [],
            expanded: true
        });
    });
    
    // 親子関係を構築
    nodeMap.forEach(node => {
        if (node.parentId && nodeMap.has(node.parentId)) {
            nodeMap.get(node.parentId).children.push(node);
        } else if (!node.parentId) {
            roots.push(node);
        }
    });
    
    return roots;
}

function renderHierarchy(nodes, container) {
    container.innerHTML = '';
    
    if (!nodes || nodes.length === 0) {
        container.innerHTML = '<div class="hierarchy-loading">該当する組織データがありません</div>';
        return;
    }
    
    nodes.forEach(node => {
        container.appendChild(createTreeNode(node, true));
    });
}

function createTreeNode(node, isRoot = false) {
    const nodeElement = document.createElement('div');
    nodeElement.className = 'tree-node' + (isRoot ? ' tree-root' : '');
    
    const contentElement = document.createElement('div');
    contentElement.className = 'tree-node-content';
    
    // トグルボタン
    if (node.children && node.children.length > 0) {
        const toggle = document.createElement('button');
        toggle.className = 'tree-toggle';
        toggle.innerHTML = node.expanded ? '▼' : '▶';
        toggle.onclick = () => toggleNode(node, nodeElement);
        contentElement.appendChild(toggle);
    } else {
        const spacer = document.createElement('span');
        spacer.className = 'tree-icon';
        spacer.innerHTML = '•';
        contentElement.appendChild(spacer);
    }
    
    // ラベル
    const label = document.createElement('span');
    label.className = 'tree-label';
    label.textContent = `${node.name} (${node.id})`;
    contentElement.appendChild(label);
    
    // 情報
    const info = document.createElement('span');
    info.className = 'tree-info';
    info.textContent = `発効: ${node.effectiveDate}`;
    if (node.expirationDate) {
        info.textContent += ` | 失効: ${node.expirationDate}`;
    }
    contentElement.appendChild(info);
    
    nodeElement.appendChild(contentElement);
    
    // 子ノード
    if (node.children && node.children.length > 0) {
        const childrenElement = document.createElement('div');
        childrenElement.className = 'tree-children' + (node.expanded ? ' expanded' : '');
        
        node.children.forEach(child => {
            childrenElement.appendChild(createTreeNode(child));
        });
        
        nodeElement.appendChild(childrenElement);
    }
    
    return nodeElement;
}

function toggleNode(node, nodeElement) {
    node.expanded = !node.expanded;
    const toggle = nodeElement.querySelector('.tree-toggle');
    const children = nodeElement.querySelector('.tree-children');
    
    if (toggle) {
        toggle.innerHTML = node.expanded ? '▼' : '▶';
    }
    
    if (children) {
        children.classList.toggle('expanded');
    }
}

function expandAll() {
    document.querySelectorAll('.tree-children').forEach(element => {
        element.classList.add('expanded');
    });
    document.querySelectorAll('.tree-toggle').forEach(button => {
        button.innerHTML = '▼';
    });
    
    // データモデルも更新
    if (hierarchyData) {
        updateExpanded(hierarchyData, true);
    }
}

function collapseAll() {
    document.querySelectorAll('.tree-children').forEach(element => {
        element.classList.remove('expanded');
    });
    document.querySelectorAll('.tree-toggle').forEach(button => {
        button.innerHTML = '▶';
    });
    
    // データモデルも更新
    if (hierarchyData) {
        updateExpanded(hierarchyData, false);
    }
}

function updateExpanded(nodes, expanded) {
    nodes.forEach(node => {
        node.expanded = expanded;
        if (node.children) {
            updateExpanded(node.children, expanded);
        }
    });
}

function exportHierarchy() {
    if (!hierarchyData) {
        alert('エクスポートする階層データがありません');
        return;
    }
    
    const date = document.getElementById('hierarchyDate').value;
    const exportData = {
        exportDate: new Date().toISOString(),
        hierarchyDate: date,
        departments: departments,
        hierarchy: hierarchyData
    };
    
    const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `organization-hierarchy-${date}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

// 組織再編機能
let reorganizeHierarchyData = null;

function updateMoveSelects() {
    const moveDeptId = document.getElementById('moveDeptId');
    const moveTargetId = document.getElementById('moveTargetId');
    
    const deptOptions = departments.map(dept => 
        `<option value="${dept.department_id}">${dept.department_id} - ${dept.department_name}</option>`
    ).join('');
    
    moveDeptId.innerHTML = '<option value="">選択してください</option>' + deptOptions;
    moveTargetId.innerHTML = '<option value="">なし（最上位へ移動）</option>' + deptOptions;
}

async function updateReorganizeView() {
    const date = document.getElementById('reorganizeDate').value;
    if (!date) return;
    
    const container = document.getElementById('reorganizeView');
    container.innerHTML = '<div class="hierarchy-loading">データを読み込み中...</div>';
    
    try {
        const response = await fetch(`${API_BASE}/hierarchy?date=${date}`);
        const data = await response.json();
        
        reorganizeHierarchyData = buildHierarchyTree(data.attributes);
        renderHierarchy(reorganizeHierarchyData, container);
    } catch (error) {
        console.error('階層データの読み込みに失敗しました:', error);
        container.innerHTML = '<div class="error">階層データの読み込みに失敗しました</div>';
    }
}

async function updateMovePreview() {
    const moveDeptId = document.getElementById('moveDeptId').value;
    const moveTargetId = document.getElementById('moveTargetId').value;
    const moveEffectiveDate = document.getElementById('moveEffectiveDate').value;
    const previewDiv = document.getElementById('movePreview');
    
    if (!moveDeptId || !moveEffectiveDate) {
        previewDiv.innerHTML = '<p>移動内容を選択してください</p>';
        return;
    }
    
    const dept = departments.find(d => d.department_id === moveDeptId);
    const targetDept = moveTargetId ? departments.find(d => d.department_id === moveTargetId) : null;
    
    // 現在の親部門を取得
    const currentDate = document.getElementById('reorganizeDate').value;
    try {
        const response = await fetch(`${API_BASE}/hierarchy?date=${currentDate}`);
        const data = await response.json();
        
        const currentAttr = data.attributes.find(a => a.department_id === moveDeptId);
        const currentParent = currentAttr && currentAttr.parent_department_id ? 
            departments.find(d => d.department_id === currentAttr.parent_department_id) : null;
        
        let previewHtml = '<h4>移動内容のプレビュー</h4>';
        previewHtml += `<p><strong>移動する部門:</strong> ${dept.department_name} (${dept.department_id})</p>`;
        previewHtml += `<p><strong>現在の上位部門:</strong> ${currentParent ? currentParent.department_name + ' (' + currentParent.department_id + ')' : 'なし（最上位）'}</p>`;
        previewHtml += `<p><strong>新しい上位部門:</strong> ${targetDept ? targetDept.department_name + ' (' + targetDept.department_id + ')' : 'なし（最上位）'}</p>`;
        previewHtml += `<p><strong>発効日:</strong> ${moveEffectiveDate}</p>`;
        
        // 循環参照チェック
        if (moveTargetId && isCircularReference(moveDeptId, moveTargetId, data.attributes)) {
            previewHtml += '<div class="preview-warning">警告: この移動により循環参照が発生します。別の部門を選択してください。</div>';
        } else if (currentAttr && currentAttr.parent_department_id === moveTargetId) {
            previewHtml += '<div class="preview-warning">情報: 現在の親部門と同じです。</div>';
        } else {
            previewHtml += '<div class="preview-success">この組織移動は実行可能です。</div>';
        }
        
        previewDiv.innerHTML = previewHtml;
    } catch (error) {
        previewDiv.innerHTML = '<div class="error">プレビューの生成に失敗しました</div>';
    }
}

function isCircularReference(moveDeptId, targetDeptId, attributes) {
    if (moveDeptId === targetDeptId) return true;
    
    const nodeMap = new Map();
    attributes.forEach(attr => {
        if (!nodeMap.has(attr.department_id)) {
            nodeMap.set(attr.department_id, []);
        }
        if (attr.parent_department_id === moveDeptId) {
            nodeMap.get(attr.department_id).push(attr.department_id);
        }
    });
    
    // 移動する部門の子孫をすべて取得
    const descendants = new Set();
    const queue = [moveDeptId];
    
    while (queue.length > 0) {
        const current = queue.shift();
        attributes.forEach(attr => {
            if (attr.parent_department_id === current && !descendants.has(attr.department_id)) {
                descendants.add(attr.department_id);
                queue.push(attr.department_id);
            }
        });
    }
    
    return descendants.has(targetDeptId);
}

async function moveOrganization(event) {
    event.preventDefault();
    
    const moveDeptId = document.getElementById('moveDeptId').value;
    const moveTargetId = document.getElementById('moveTargetId').value || null;
    const moveEffectiveDate = document.getElementById('moveEffectiveDate').value;
    
    if (!moveDeptId || !moveEffectiveDate) {
        alert('必要な項目を入力してください');
        return;
    }
    
    // 循環参照チェック
    const currentDate = document.getElementById('reorganizeDate').value;
    try {
        const response = await fetch(`${API_BASE}/hierarchy?date=${currentDate}`);
        const data = await response.json();
        
        if (moveTargetId && isCircularReference(moveDeptId, moveTargetId, data.attributes)) {
            alert('この移動により循環参照が発生します。別の部門を選択してください。');
            return;
        }
    } catch (error) {
        alert('現在の組織構造の確認に失敗しました');
        return;
    }
    
    // 新しい組織属性レコードを作成
    const newAttribute = {
        department_id: moveDeptId,
        effective_date: moveEffectiveDate,
        parent_department_id: moveTargetId
    };
    
    try {
        const response = await fetch(`${API_BASE}/organization-attributes`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newAttribute)
        });
        
        if (response.ok) {
            alert('組織移動が正常に登録されました');
            loadAttributes();
            updateReorganizeView();
            
            // フォームをリセット
            document.getElementById('moveDeptId').value = '';
            document.getElementById('moveTargetId').value = '';
            updateMovePreview();
        } else {
            const error = await response.json();
            alert('エラー: ' + error.error);
        }
    } catch (error) {
        alert('組織移動の登録に失敗しました: ' + error.message);
    }
}