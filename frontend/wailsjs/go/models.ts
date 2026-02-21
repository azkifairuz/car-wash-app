export namespace backend {
	
	export class CartItem {
	    itemType: string;
	    itemId: number;
	    itemName: string;
	    price: number;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new CartItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemType = source["itemType"];
	        this.itemId = source["itemId"];
	        this.itemName = source["itemName"];
	        this.price = source["price"];
	        this.quantity = source["quantity"];
	    }
	}

}

export namespace models {
	
	export class DashboardStats {
	    todayRevenue: number;
	    todayTransactions: number;
	    activeWashes: number;
	    foodOrdersToday: number;
	    monthlyRevenue: number;
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.todayRevenue = source["todayRevenue"];
	        this.todayTransactions = source["todayTransactions"];
	        this.activeWashes = source["activeWashes"];
	        this.foodOrdersToday = source["foodOrdersToday"];
	        this.monthlyRevenue = source["monthlyRevenue"];
	    }
	}
	export class Discount {
	    id: number;
	    name: string;
	    type: string;
	    value: number;
	    minOrder: number;
	    isActive: boolean;
	    validFrom: string;
	    validUntil: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Discount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.value = source["value"];
	        this.minOrder = source["minOrder"];
	        this.isActive = source["isActive"];
	        this.validFrom = source["validFrom"];
	        this.validUntil = source["validUntil"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class User {
	    id: number;
	    username: string;
	    fullName: string;
	    role: string;
	    isActive: boolean;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.fullName = source["fullName"];
	        this.role = source["role"];
	        this.isActive = source["isActive"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class LoginResponse {
	    success: boolean;
	    message: string;
	    user?: User;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.user = this.convertValues(source["user"], User);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MenuItem {
	    id: number;
	    name: string;
	    price: number;
	    category: string;
	    icon: string;
	    stock: number;
	    isActive: boolean;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new MenuItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.price = source["price"];
	        this.category = source["category"];
	        this.icon = source["icon"];
	        this.stock = source["stock"];
	        this.isActive = source["isActive"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class Shift {
	    id: number;
	    userId: number;
	    user: User;
	    date: string;
	    startTime: string;
	    endTime: string;
	    notes: string;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Shift(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.userId = source["userId"];
	        this.user = this.convertValues(source["user"], User);
	        this.date = source["date"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TransactionItem {
	    id: number;
	    transactionId: number;
	    itemType: string;
	    itemId: number;
	    itemName: string;
	    price: number;
	    quantity: number;
	    subtotal: number;
	
	    static createFrom(source: any = {}) {
	        return new TransactionItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.transactionId = source["transactionId"];
	        this.itemType = source["itemType"];
	        this.itemId = source["itemId"];
	        this.itemName = source["itemName"];
	        this.price = source["price"];
	        this.quantity = source["quantity"];
	        this.subtotal = source["subtotal"];
	    }
	}
	export class Transaction {
	    id: number;
	    invoiceNo: string;
	    customerName: string;
	    plateNumber: string;
	    carType: string;
	    items: TransactionItem[];
	    subtotal: number;
	    discountId?: number;
	    discount?: Discount;
	    discountAmount: number;
	    total: number;
	    paymentMethod: string;
	    status: string;
	    cashierId: number;
	    cashier: User;
	    createdAt: string;
	    updatedAt: string;
	    completedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Transaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoiceNo = source["invoiceNo"];
	        this.customerName = source["customerName"];
	        this.plateNumber = source["plateNumber"];
	        this.carType = source["carType"];
	        this.items = this.convertValues(source["items"], TransactionItem);
	        this.subtotal = source["subtotal"];
	        this.discountId = source["discountId"];
	        this.discount = this.convertValues(source["discount"], Discount);
	        this.discountAmount = source["discountAmount"];
	        this.total = source["total"];
	        this.paymentMethod = source["paymentMethod"];
	        this.status = source["status"];
	        this.cashierId = source["cashierId"];
	        this.cashier = this.convertValues(source["cashier"], User);
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.completedAt = source["completedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class WashPackage {
	    id: number;
	    name: string;
	    description: string;
	    price: number;
	    duration: number;
	    category: string;
	    icon: string;
	    isActive: boolean;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new WashPackage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.price = source["price"];
	        this.duration = source["duration"];
	        this.category = source["category"];
	        this.icon = source["icon"];
	        this.isActive = source["isActive"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}

}

