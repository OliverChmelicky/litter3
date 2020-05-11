import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditTrashComponent } from './edit-trash.component';

describe('EditTrashComponent', () => {
  let component: EditTrashComponent;
  let fixture: ComponentFixture<EditTrashComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EditTrashComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditTrashComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
