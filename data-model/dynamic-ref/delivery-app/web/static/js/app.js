let currentTab = 'customers';
let prefecturesData = [];
let regionDeliveryMethods = [];

async function loadData() {
    await loadPrefectures();
    await loadRegionDeliveryMethods();
    await loadCustomers();
}

async function loadPrefectures() {
    try {
        const response = await fetch('/api/prefectures');
        prefecturesData = await response.json();
        
        const select = document.getElementById('prefecture');
        select.innerHTML = '<option value="">選択してください</option>';
        prefecturesData.forEach(pref => {
            const option = document.createElement('option');
            option.value = pref.prefecture;
            option.textContent = pref.prefecture;
            select.appendChild(option);
        });
        
        updatePrefecturesTable();
    } catch (error) {
        console.error('Error loading prefectures:', error);
    }
}

async function loadRegionDeliveryMethods() {
    try {
        const response = await fetch('/api/region-delivery-methods');
        regionDeliveryMethods = await response.json();
        updateRegionDeliveryMethodsTable();
    } catch (error) {
        console.error('Error loading region delivery methods:', error);
    }
}

async function loadCustomers() {
    try {
        const response = await fetch('/api/customers');
        const customers = await response.json();
        updateCustomersTable(customers);
    } catch (error) {
        console.error('Error loading customers:', error);
    }
}

function showTab(tabName) {
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab-button').forEach(button => {
        button.classList.remove('active');
    });
    
    document.getElementById(tabName).classList.add('active');
    event.target.classList.add('active');
    currentTab = tabName;
}

function showCustomerForm() {
    document.getElementById('customerForm').style.display = 'block';
    document.getElementById('customerId').value = '';
    document.getElementById('customerName').value = '';
    document.getElementById('prefecture').value = '';
    document.getElementById('deliveryMethod').innerHTML = '<option value="">都道府県を選択してください</option>';
}

function hideCustomerForm() {
    document.getElementById('customerForm').style.display = 'none';
}

function updateDeliveryMethods() {
    const prefecture = document.getElementById('prefecture').value;
    const deliveryMethodSelect = document.getElementById('deliveryMethod');
    
    if (!prefecture) {
        deliveryMethodSelect.innerHTML = '<option value="">都道府県を選択してください</option>';
        return;
    }
    
    const prefData = prefecturesData.find(p => p.prefecture === prefecture);
    if (!prefData) return;
    
    const availableMethods = regionDeliveryMethods.filter(m => m.region === prefData.region);
    
    deliveryMethodSelect.innerHTML = '<option value="">選択してください</option>';
    availableMethods.forEach(method => {
        const option = document.createElement('option');
        option.value = method.delivery_method;
        option.textContent = `${method.delivery_method} (${method.standard_delivery_days}日, ¥${method.standard_delivery_fee})`;
        deliveryMethodSelect.appendChild(option);
    });
}

async function saveCustomer(event) {
    event.preventDefault();
    
    const customerId = document.getElementById('customerId').value;
    const customerData = {
        customer_name: document.getElementById('customerName').value,
        prefecture: document.getElementById('prefecture').value,
        delivery_method: document.getElementById('deliveryMethod').value
    };
    
    try {
        let response;
        if (customerId) {
            response = await fetch(`/api/customers/${customerId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(customerData)
            });
        } else {
            response = await fetch('/api/customers', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(customerData)
            });
        }
        
        const result = await response.json();
        if (response.ok) {
            hideCustomerForm();
            loadCustomers();
        } else {
            alert(result.error || 'エラーが発生しました');
        }
    } catch (error) {
        alert('保存中にエラーが発生しました');
    }
}

function editCustomer(customer) {
    showCustomerForm();
    document.getElementById('customerId').value = customer.customer_id;
    document.getElementById('customerName').value = customer.customer_name;
    document.getElementById('prefecture').value = customer.prefecture;
    updateDeliveryMethods();
    setTimeout(() => {
        document.getElementById('deliveryMethod').value = customer.delivery_method;
    }, 100);
}

function updateCustomersTable(customers) {
    const tbody = document.querySelector('#customersTable tbody');
    tbody.innerHTML = '';
    
    customers.forEach(customer => {
        const row = tbody.insertRow();
        row.innerHTML = `
            <td>${customer.customer_id}</td>
            <td>${customer.customer_name}</td>
            <td>${customer.prefecture}</td>
            <td>${customer.region}</td>
            <td>${customer.delivery_method}</td>
            <td>${customer.standard_delivery_days}日</td>
            <td>¥${customer.standard_delivery_fee}</td>
            <td>
                <button class="btn-edit" onclick='editCustomer(${JSON.stringify(customer)})'>編集</button>
            </td>
        `;
    });
}

function updateRegionDeliveryMethodsTable() {
    const tbody = document.querySelector('#regionsTable tbody');
    tbody.innerHTML = '';
    
    regionDeliveryMethods.forEach(method => {
        const row = tbody.insertRow();
        row.innerHTML = `
            <td>${method.region}</td>
            <td>${method.delivery_method}</td>
            <td>${method.standard_delivery_days}日</td>
            <td>¥${method.standard_delivery_fee}</td>
            <td>
                <button class="btn-danger" onclick="deleteRegionDeliveryMethod('${method.region}', '${method.delivery_method}')">削除</button>
            </td>
        `;
    });
}

async function deleteRegionDeliveryMethod(region, deliveryMethod) {
    if (!confirm(`${region}の${deliveryMethod}を削除しますか？`)) return;
    
    try {
        const response = await fetch(`/api/region-delivery-methods?region=${region}&delivery_method=${deliveryMethod}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        if (response.ok) {
            loadRegionDeliveryMethods();
        } else {
            alert(result.error || '削除に失敗しました');
        }
    } catch (error) {
        alert('削除中にエラーが発生しました');
    }
}

function updatePrefecturesTable() {
    const tbody = document.querySelector('#prefecturesTable tbody');
    tbody.innerHTML = '';
    
    prefecturesData.forEach(pref => {
        const row = tbody.insertRow();
        row.innerHTML = `
            <td>${pref.prefecture}</td>
            <td>
                <select onchange="updatePrefectureRegion('${pref.prefecture}', this.value)">
                    <option value="関西" ${pref.region === '関西' ? 'selected' : ''}>関西</option>
                    <option value="中部" ${pref.region === '中部' ? 'selected' : ''}>中部</option>
                </select>
            </td>
            <td>
                <button class="btn-primary" onclick="updatePrefectureRegion('${pref.prefecture}', this.previousElementSibling.children[0].value)">更新</button>
            </td>
        `;
    });
}

async function updatePrefectureRegion(prefecture, newRegion) {
    try {
        const response = await fetch(`/api/prefectures/${prefecture}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ region: newRegion })
        });
        
        const result = await response.json();
        if (response.ok) {
            loadPrefectures();
            loadCustomers();
        } else {
            alert(result.error || '更新に失敗しました');
        }
    } catch (error) {
        alert('更新中にエラーが発生しました');
    }
}

window.onload = loadData;