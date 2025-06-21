import { defineStore } from 'pinia'

export const useMessageStore = defineStore('message', {
  state: () => ({
    contacts: [],
    messages: {},
    currentContactId: null,
    loadingContacts: false,
    loadingMessages: false,
    totalUnread: 0,
    // 添加联系人商品关联映射 { contactId: productId }
    contactProductMap: {},
    // 添加联系人多商品映射 { contactId: [productIds] }
    contactMultiProductsMap: {},
    contactList: []
  }),
  
  getters: {
    // 获取当前联系人
    currentContact: (state) => {
      if (!state.currentContactId) return null;
      return state.contacts.find(c => c.id === state.currentContactId) || null;
    },
    
    // 获取当前联系人的消息
    currentMessages: (state) => {
      if (!state.currentContactId) return [];
      
      // 先从contacts中查找
      const contact = state.contacts.find(c => c.id === state.currentContactId);
      if (contact && contact.messages) {
        return contact.messages;
      }
      
      // 后备方案：从messages对象中获取
      return state.messages[state.currentContactId] || [];
    },
    
    // 根据ID获取联系人
    getContactById: (state) => (id) => {
      // 如果id是字符串，转换为数字
      const contactId = typeof id === 'string' ? parseInt(id) : id;
      return state.contacts.find(c => c.id === contactId) || null;
    },
    
    // 过滤联系人
    filterContacts: (state) => (keyword) => {
      if (!keyword) return state.contacts;
      
      return state.contacts.filter(contact => 
        contact.username && contact.username.toLowerCase().includes(keyword.toLowerCase())
      );
    },
    
    // 获取联系人关联的商品
    getContactProduct: (state) => (contactId) => {
      return state.contactProductMap[contactId] || null;
    },
    
    // 获取当前联系人关联的商品
    currentContactProduct: (state) => {
      if (!state.currentContactId) return null;
      return state.contactProductMap[state.currentContactId] || null;
    },
    
    // 获取联系人关联的多个商品
    getContactProducts: (state) => (contactId) => {
      return state.contactMultiProductsMap[contactId] || [];
    },
    
    // 获取当前联系人关联的多个商品
    currentContactProducts: (state) => {
      if (!state.currentContactId) return [];
      return state.contactMultiProductsMap[state.currentContactId] || [];
    }
  },
  
  actions: {
    // 设置联系人列表
    setContacts(contacts) {
      this.contacts = contacts;
      
      // 计算总未读数量
      this.updateTotalUnread();
    },
    
    // 添加或更新联系人
    addOrUpdateContact(contact) {
      const index = this.contacts.findIndex(c => c.id === contact.id);
      
      if (index > -1) {
        // 更新现有联系人
        this.contacts[index] = { ...this.contacts[index], ...contact };
      } else {
        // 添加新联系人到列表顶部
        this.contacts.unshift(contact);
      }
      
      // 如果联系人有商品ID且大于0，保存到映射
      if (contact.product_id && contact.product_id > 0) {
        this.setContactProduct(contact.id, contact.product_id);
      }
      
      // 更新总未读数量
      this.updateTotalUnread();
    },
    
    // 设置指定联系人的消息
    setMessages(contactId, messages) {
      this.messages[contactId] = messages;
    },
    
    // 添加消息 - 重写以提高稳定性和可靠性
    addMessage(contactId, messageData) {
      if (!contactId || !messageData) {
        console.warn('添加消息失败: 缺少必要参数', { contactId, messageData });
        return;
      }
      
      try {
        // 确保contactId是数字
        const numericContactId = typeof contactId === 'string' ? parseInt(contactId) : contactId;
        
        // 格式化消息数据，确保必要字段存在
        const message = {
          id: messageData.id || Date.now(),
          content: messageData.content || '',
          sender_id: messageData.sender_id || 0,
          receiver_id: messageData.receiver_id || 0,
          created_at: messageData.created_at || new Date().toISOString(),
          is_read: messageData.is_read || false
        };
        
        // 如果存在product_id且有效，则添加
        if (messageData.product_id && messageData.product_id > 0) {
          message.product_id = messageData.product_id;
        }
        
        console.log('准备添加消息到store:', { contactId: numericContactId, message });
        
        // 确保messages对象中有联系人消息数组
        if (!this.messages[numericContactId]) {
          this.messages[numericContactId] = [];
        }
        
        // 检查消息是否已存在
        const existingMsgIndex = this.messages[numericContactId].findIndex(m => m.id === message.id);
        
        if (existingMsgIndex >= 0) {
          // 更新已存在的消息
          this.messages[numericContactId][existingMsgIndex] = {
            ...this.messages[numericContactId][existingMsgIndex],
            ...message
          };
        } else {
          // 添加新消息
          this.messages[numericContactId].push(message);
          
          // 确保消息按时间顺序排列
          this.messages[numericContactId].sort((a, b) => {
            const timeA = new Date(a.created_at || a.time).getTime();
            const timeB = new Date(b.created_at || b.time).getTime();
            return timeA - timeB;
          });
        }
        
        // 查找或创建联系人
        let contact = this.contacts.find(c => c.id === numericContactId);
        
        // 如果联系人不存在，创建新联系人
        if (!contact) {
          console.log('联系人不存在，创建新联系人:', numericContactId);
          contact = {
            id: numericContactId,
            username: `联系人${numericContactId}`,
            avatar: '',
            lastMessage: message.content,
            lastActiveTime: message.created_at,
            unread: 0
          };
          this.contacts.push(contact);
        }
        
        // 更新联系人最后消息和时间
        contact.lastMessage = message.content;
        contact.lastActiveTime = message.created_at;
        
        // 如果消息带有商品ID，更新商品关联
        if (message.product_id && message.product_id > 0) {
          this.addContactProduct(numericContactId, message.product_id);
        }
        
        // 更新联系人顺序
        this.updateContactsOrder();
        
        console.log('消息添加成功，更新后的联系人:', contact);
      } catch (error) {
        console.error('添加消息过程中出错:', error);
      }
    },
    
    // 设置联系人关联的商品ID
    setContactProduct(contactId, productId) {
      if (!contactId || !productId || productId <= 0) {
        delete this.contactProductMap[contactId];
      } else {
        this.contactProductMap[contactId] = productId;
        // 持久化到localStorage
        this.saveProductMappings();
      }
    },
    
    // 添加商品ID到联系人的多商品列表
    addContactProduct(contactId, productId) {
      if (!contactId || !productId || productId <= 0) return;
      
      if (!this.contactMultiProductsMap[contactId]) {
        this.contactMultiProductsMap[contactId] = [];
      }
      
      // 如果商品ID不在列表中，添加到列表
      if (!this.contactMultiProductsMap[contactId].includes(productId)) {
        this.contactMultiProductsMap[contactId].push(productId);
        // 持久化到localStorage
        this.saveProductMappings();
      }
    },
    
    // 从联系人的多商品列表中移除商品ID
    removeContactProduct(contactId, productId) {
      if (!this.contactMultiProductsMap[contactId]) return;
      
      const index = this.contactMultiProductsMap[contactId].indexOf(productId);
      if (index > -1) {
        this.contactMultiProductsMap[contactId].splice(index, 1);
        // 持久化到localStorage
        this.saveProductMappings();
      }
    },
    
    // 持久化商品映射到localStorage
    saveProductMappings() {
      try {
        localStorage.setItem('contactProductMap', JSON.stringify(this.contactProductMap));
        localStorage.setItem('contactMultiProductsMap', JSON.stringify(this.contactMultiProductsMap));
      } catch (e) {
        console.error('保存商品映射到localStorage失败:', e);
      }
    },
    
    // 从localStorage恢复商品映射
    restoreProductMappings() {
      try {
        const savedProductMap = localStorage.getItem('contactProductMap');
        const savedMultiProductsMap = localStorage.getItem('contactMultiProductsMap');
        
        if (savedProductMap) {
          this.contactProductMap = JSON.parse(savedProductMap);
        }
        
        if (savedMultiProductsMap) {
          this.contactMultiProductsMap = JSON.parse(savedMultiProductsMap);
        }
        
        console.log('已从localStorage恢复商品映射');
      } catch (e) {
        console.error('从localStorage恢复商品映射失败:', e);
      }
    },
    
    // 设置当前联系人
    setCurrentContact(contactId) {
      this.currentContactId = contactId;
      
      // 清除未读标记
      if (contactId) {
        const contact = this.contacts.find(c => c.id === contactId);
        if (contact && contact.unread) {
          contact.unread = 0;
          
          // 更新总未读数量
          this.updateTotalUnread();
        }
      }
    },
    
    // 清除指定联系人的未读数量
    clearUnread(contactId) {
      const contact = this.contacts.find(c => c.id === contactId);
      if (contact && contact.unread) {
        contact.unread = 0;
        
        // 更新总未读数量
        this.updateTotalUnread();
      }
    },
    
    // 更新总未读数量
    updateTotalUnread() {
      if (!this.contacts || !Array.isArray(this.contacts)) {
        console.warn('无法更新未读消息数量：contacts不是数组', this.contacts);
        this.totalUnread = 0;
        return;
      }
      
      this.totalUnread = this.contacts.reduce((total, contact) => {
        return total + (contact.unread || 0);
      }, 0);
    },
    
    // 将联系人移动到列表顶部
    moveContactToTop(contactId) {
      const index = this.contacts.findIndex(c => c.id === contactId);
      if (index > 0) {
        // 从当前位置删除
        const contact = this.contacts.splice(index, 1)[0];
        // 添加到顶部
        this.contacts.unshift(contact);
      }
    },
    
    // 设置加载状态
    setLoading(type, status) {
      if (type === 'contacts') {
        this.loadingContacts = status;
      } else if (type === 'messages') {
        this.loadingMessages = status;
      }
    },
    
    // 清空所有数据
    clearAll() {
      this.contacts = [];
      this.messages = {};
      this.currentContactId = null;
      this.totalUnread = 0;
      this.loadingContacts = false;
      this.loadingMessages = false;
      // 不清除商品映射，因为用户可能会再次登录
    },
    
    // Add new method
    updateContactWithMessage(contactId, message) {
      const contact = this.contacts.find(c => c.id === contactId);
      
      if (contact) {
        // Update existing contact
        contact.lastMessage = message.content;
        contact.lastTime = message.created_at || message.time;
        
        if (contactId !== this.currentContactId) {
          contact.unread = (contact.unread || 0) + 1;
          this.updateTotalUnread();
        }
        
        this.moveContactToTop(contactId);
      } else {
        // We need to fetch contact info and create a new contact
        console.warn(`Contact with ID ${contactId} not found, need to fetch info`);
      }
    },
    
    setTotalUnread(count) {
      this.totalUnread = count;
    },
    
    // 更新联系人顺序，将有最新消息的联系人移到最前
    updateContactsOrder() {
      try {
        if (!this.contacts || this.contacts.length === 0) {
          return; // 防止空数组操作
        }
        
        // 复制数组防止直接修改状态
        const sortedContacts = [...this.contacts];
        
        // 根据最后活跃时间排序
        sortedContacts.sort((a, b) => {
          // 安全获取时间，提供默认值防止错误
          const timeA = a && a.lastActiveTime ? new Date(a.lastActiveTime).getTime() : 0;
          const timeB = b && b.lastActiveTime ? new Date(b.lastActiveTime).getTime() : 0;
          
          // 降序排列，最新的在前
          return timeB - timeA;
        });
        
        // 更新contacts数组
        this.contacts = sortedContacts;
      } catch (error) {
        console.error('排序联系人时出错:', error);
        // 发生错误时不改变现有顺序
      }
    }
  }
}) 