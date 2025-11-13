import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class AuthService {
  auth: { first_name: string } | null = null;

  login(name: string) {
    this.auth = { first_name: name };
  }

  logout() {
    this.auth = null;
  }
}
