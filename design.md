## 3. **Functional Requirements**

### 3.1 Book Management

- [x] FR1: The system shall allow librarians to add new books with details (title, author, ISBN, category, number of copies).
- [x] FR5: The system shall track both total and available copies of each book.
- [x] FR3: The system shall allow deleting books only if there are no active loans associated with them.
- [x] FR4: The system shall allow searching for books by title, author, ISBN, or category.
- [x] FR2: The system shall allow updating book details (title, author, category, etc.).

### 3.2 Member Management

- [x] FR6: The system shall allow adding new members with details (name, email, phone).
    
- [x] FR7: The system shall allow updating member details.
    
- [ ] FR8: The system shall prevent deletion of a member if they have outstanding loans.
    
- [x] FR9: The system shall allow searching for members by name or ID.
    

### 3.3 Borrowing and Returning

- [ ] The system shall allow marking loans as returned

- [ ] FR10: The system shall allow borrowing books only if at least one copy is available.
    
- [x] FR11: The system shall create a loan record for each borrowed book, including borrow date and due date (default 14 days).
    
- [ ] FR12: The system shall allow returning books, updating the loan record with a return date and increasing the book’s available copies.
    
- [ ] FR13: The system shall mark loans as “overdue” if the current date exceeds the due date.
    
- FR14 (Optional): The system shall calculate late fees for overdue books.
    

### 3.4 Loan Tracking and Reports

- FR15: The system shall allow librarians to list all active loans.
    
- FR16: The system shall allow viewing all overdue loans.
    
- FR17: The system shall generate a report of most borrowed books.
    
- FR18: The system shall generate a report of active members and their loan history.
