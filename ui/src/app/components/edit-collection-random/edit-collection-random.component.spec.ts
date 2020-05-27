import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditCollectionRandomComponent } from './edit-collection-random.component';

describe('EditCollectionRandomComponent', () => {
  let component: EditCollectionRandomComponent;
  let fixture: ComponentFixture<EditCollectionRandomComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EditCollectionRandomComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditCollectionRandomComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
